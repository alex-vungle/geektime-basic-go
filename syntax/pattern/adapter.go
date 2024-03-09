package pattern

import "context"

type MyServiceV1 interface {
	Handle(req any)
}

type MyServiceV2 interface {
	Handle(ctx context.Context, req any)
}

type v2Tov1 struct {
	v2 MyServiceV2
}

func (v *v2Tov1) Handle(req any) {
	v.v2.Handle(context.Background(), req)
}

type MyServiceInvoker struct {
	v1 MyServiceV1
}

func main() {
	var v2 MyServiceV2
	NewMyServiceInvoker(&v2Tov1{
		v2: v2,
	})
}

func NewMyServiceInvoker(v1 MyServiceV1) MyServiceInvoker {
	return MyServiceInvoker{
		v1: v1,
	}
}
