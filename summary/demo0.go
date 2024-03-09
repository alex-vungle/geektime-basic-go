package summary

import (
	"context"
	"gitee.com/geekbang/basic-go/summary/somefunc"
)

type YourBusiness struct {
	doer somefunc.DoSomething
}

func (b *YourBusiness) DoBiz() {
	// 最好不要
	// 你没办法 mock，你也没办法换实现
	somefunc.Do(context.Background(), 123)
	// 你应该调这个
	//b.doer.Do(context.Background(), 123)
}

var handlers map[int]func()

func DemoManySwitchCases(enum int) {
	switch enum {
	case 123:
	case 124:
	case 125:
	case 126:

	}
}

// Enum 把 switch case 变成不同的实现
type Enum interface {
}

type Filter interface {
	DoFilter(val int) error
}

type myFilter struct {
}

//func (f *myFilter)DoFilter(val int) error {
//	if (val < 10) {
//		// xxxx
//	}
//	// 不然，下一步
//}
