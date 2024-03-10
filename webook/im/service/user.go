package service

import (
	"context"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/im/domain"
	"github.com/ecodeclub/ekit/net/httpx"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type UserService interface {
	Sync(ctx context.Context, user domain.User) error
	//SyncUsers(ctx context.Context, user []domain.User) error
}

type RESTUserService struct {
	// HTTP 请求的域名端口
	base string
	// 默认是 openIM123
	secret string
	// 一旦将来你要换 client，你很容易就换掉
	client *http.Client
}

func NewRESTUserService(base string, secret string) *RESTUserService {
	return &RESTUserService{base: base, secret: secret,
		client: http.DefaultClient}
}

func (svc *RESTUserService) Sync(ctx context.Context, user domain.User) error {
	var operationID string
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		operationID = spanCtx.TraceID().String()
	} else {
		operationID = uuid.New().String()
	}
	var resp response
	err := httpx.NewRequest(ctx,
		http.MethodPost,
		svc.base+"/user/user_register").
		AddHeader("operationID", operationID).JSONBody(request{
		Secret: svc.secret,
		Users:  []domain.User{user},
	}).Client(svc.client).Do().JSONScan(&resp)
	if err != nil {
		return err
	}
	if resp.ErrCode != 0 {
		return fmt.Errorf("同步用户数据失败 %v", resp)
	}
	return nil
}

type request struct {
	Secret string        `json:"secret"`
	Users  []domain.User `json:"users"`
}

type response struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
}
