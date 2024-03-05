package tencent

import (
	"context"
	"fmt"
	sms2 "gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"go.uber.org/zap"
)

type Service struct {
	client   *sms.Client
	appId    *string
	signName *string
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = s.appId
	request.SignName = s.signName
	request.TemplateId = ekit.ToPtr[string](tplId)
	request.TemplateParamSet = s.toPtrSlice(args)
	request.PhoneNumberSet = s.toPtrSlice(numbers)
	response, err := s.client.SendSms(request)
	zap.L().Debug("请求腾讯SendSMS接口",
		zap.Any("req", request),
		zap.Any("resp", response))
	// 处理异常
	if err != nil {
		return err
	}
	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr == nil {
			// 不可能进来这里
			continue
		}
		status := *statusPtr
		if status.Code == nil || *(status.Code) != "Ok" {
			// 发送失败
			return fmt.Errorf("发送短信失败 code: %s, msg: %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) toPtrSlice(data []string) []*string {
	return slice.Map[string, *string](data,
		func(idx int, src string) *string {
			return &src
		})
}

func NewService(client *sms.Client, appId string, signName string) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		signName: &signName,
	}
}

func (s *Service) SendV11(ctx context.Context, tplId string, args any, numbers ...string) error {
	return s.Send(ctx, tplId, args.([]string), numbers...)
}

func (s *Service) SendV22(ctx context.Context, tplId string, args []sms2.NamedArg, numbers ...string) error {
	// 转切片
	newArgs := slice.Map(args, func(idx int, src sms2.NamedArg) string {
		return src.Value
	})
	return s.Send(ctx, tplId, newArgs, numbers...)
}
