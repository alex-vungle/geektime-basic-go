package service

import (
	"context"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/im/domain"
	"github.com/ecodeclub/ekit/net/httpx"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
)

// 你可以用 localhost
const defaultEndpoint = "http://localhost:10002/user/user_register"

type UserService interface {
	SyncUser(ctx context.Context, user domain.User) error
}

type RESTUserService struct {
	endpoint string
	secret   string
	client   *http.Client
}

// secret 从哪里来？
// 默认就是 openIM123
func NewUserService(secret string) *RESTUserService {
	// 假如说我有 TLS 之类的认证，我在这里可以灵活替换具体的 client
	return &RESTUserService{endpoint: defaultEndpoint,
		client: http.DefaultClient}
}

func (s *RESTUserService) SyncUser(ctx context.Context, user domain.User) error {
	spanCtx := trace.SpanContextFromContext(ctx)
	// 如果你本身有链路追踪
	var operationID string
	// 有这个
	if spanCtx.HasTraceID() {
		operationID = spanCtx.TraceID().String()
	} else {
		// 实际上也就是你这边没有接入 otel，或者断开了。你的链路追踪断开了
		operationID = uuid.New().String()
	}
	// 现在也就是我要发请求调用 openim 的接口了
	// httpx 是我封装过的
	// 这个东西本身就是批量接口
	var resp response

	err := httpx.NewRequest(ctx, http.MethodPost, s.endpoint).
		JSONBody(Request{Secret: s.secret, Users: []User{
			{
				UserID:   strconv.FormatInt(user.ID, 10),
				Nickname: user.Nickname,
				FaceURL:  user.Avatar,
			}}}).
		Client(s.client).
		AddHeader("operationID", operationID).
		// 发起请求，拿到 response
		Do().
		// 我把请求体转为一个结构体
		JSONScan(&resp)
	if err != nil {
		return err
	}
	if resp.ErrCode != 0 {
		// 也是出错了
		return fmt.Errorf("同步用户数据失败 %v", resp)
	}
	return nil
}

type response struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
}

type Request struct {
	Secret string `json:"secret"`
	Users  []User `json:"users"`
}

type User struct {
	UserID   string `json:"userID"`
	Nickname string `json:"nickname"`
	// 头像
	FaceURL string `json:"faceURL"`
}
