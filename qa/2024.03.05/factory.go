package qa

type Car struct {
}

type CarFactory func() Car
type CarFactoryV1 func(cfg string) Car

type CarFactoryV2 interface {
	Create() Car
}
