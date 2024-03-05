package sms

import "context"

// Service 发送短信的抽象
// 屏蔽不同供应商之间的区别
//
//go:generate mockgen -source=./types.go -package=smsmocks -destination=./mocks/sms.mock.go Service
type Service interface {
	Send(ctx context.Context, tplId string,
		args []string, numbers ...string) error
	//A()
	//B()
}

// 我在设计的时候就知道，阿里云和腾讯云的需要的东西不一样
// 直接设计一个兼容了两种形态的接口
type ServiceV1 interface {
	SendV11(ctx context.Context, tplId string,
		// args 直接 any。
		// 腾讯云解释 args 为 []string
		// 阿里云解释 args 为 map[string]string
		// 缺陷就是，codeService 那边，需要能够正确构造 args
		args any, numbers ...string) error
}

type ServiceV2 interface {
	SendV22(ctx context.Context, tplId string,
		// 缺陷就是，codeService 那边，需要能够正确构造 args
		args []NamedArg, numbers ...string) error

	SendV23(ctx context.Context, tplId string, args []NamedArg, newArgs int, numbers ...string) error
}

type ServiceFactory func() Service

type NamedArg struct {
	Name  string
	Value string
}

type Req struct {
}
