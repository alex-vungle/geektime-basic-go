package somefunc

import "context"

type DoSomething interface {
	Do(ctx context.Context, req any) error
}

var def DoSomething

func Do(ctx context.Context, req any) error {
	return def.Do(ctx, req)
}
