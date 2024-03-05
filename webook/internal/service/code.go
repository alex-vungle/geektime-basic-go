package service

import (
	"context"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"math/rand"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

//go:generate mockgen -source=./code.go -package=svcmocks -destination=./mocks/code.mock.go CodeService
type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context,
		biz, phone, inputCode string) (bool, error)
}

// 面向接口编程，绝对不会被具体实现影响
// 随便你换实现
type codeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{
		repo: repo,
		sms:  smsSvc,
	}
}

func (svc *codeService) Send(ctx context.Context, biz, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	// 你在这儿，是不是要开始发送验证码了？
	if err != nil {
		return err
	}
	const codeTplId = "1877556"
	// 假如说你没有接口，
	// 没有接口，你替换实现的时候，就很难替换阿里云，
	// 破坏了面向接口编程，依赖于实现细节（阿里云要 map，而腾讯云要切片）
	//args := map[string]string{
	//	"code": code,
	//}
	return svc.sms.Send(ctx, codeTplId, []string{code}, phone)

	// 这种形态下，你依旧保持了面向接口编程，不依赖于细节
	//return svc.sms.SendV22(ctx, codeTplId,
	//	[]sms.NamedArg{{Name: "code", Value: code}}, phone)

	//
}

func (svc *codeService) Verify(ctx context.Context,
	biz, phone, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if err == repository.ErrCodeVerifyTooMany {
		// 相当于，我们对外面屏蔽了验证次数过多的错误，我们就是告诉调用者，你这个不对
		return false, nil
	}
	return ok, err
}

func (svc *codeService) generate() string {
	// 0-999999
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
