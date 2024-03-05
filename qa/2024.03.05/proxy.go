package qa

type Service interface {
	Serve()
}

// 真正的动态代理，是运行的时候，为 Service 提供一个实现

type ServiceV1 struct {
	Serve func()
}
