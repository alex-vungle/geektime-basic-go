package wrr

import (
	"context"
	"errors"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

const Name = "custom_weighted_round_robin"

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &PickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type PickerBuilder struct {
}

func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*weightConn, 0, len(info.ReadySCs))
	for sc, sci := range info.ReadySCs {
		md, _ := sci.Address.Metadata.(map[string]any)
		weightVal, _ := md["weight"]
		weight, _ := weightVal.(float64)
		//if weight == 0 {
		//
		//}
		conns = append(conns, &weightConn{
			SubConn:       sc,
			weight:        int(weight),
			currentWeight: int(weight),
		})
	}

	return &Picker{
		conns: conns,
	}
}

type Picker struct {
	conns []*weightConn
	lock  sync.Mutex
}

func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var totalWeight int
	var res *weightConn

	for _, c := range p.conns {
		c.mutex.Lock()
		totalWeight = totalWeight + c.efficientWeight
		c.currentWeight = c.currentWeight + c.efficientWeight
		if res == nil || res.currentWeight < c.currentWeight {
			res = c
		}
		c.mutex.Unlock()
	}
	res.mutex.Lock()
	res.currentWeight = res.currentWeight - totalWeight
	res.mutex.Unlock()
	return balancer.PickResult{
		SubConn: res.SubConn,
		Done: func(info balancer.DoneInfo) {
			res.mutex.Lock()
			defer res.mutex.Unlock()
			if info.Err != nil && res.efficientWeight == 0 {
				return
			}
			switch {
			case info.Err == nil:
				// 假设权重最大为400
				if res.efficientWeight == 400 {
					return
				}
				// 增加权重
				res.efficientWeight++
			case errors.Is(info.Err, context.DeadlineExceeded):
				// 动态调整超时
				res.efficientWeight = res.efficientWeight - 10
			default:
				// 服务端错误
				code := status.Code(info.Err)
				switch code {
				case codes.Unavailable:
					res.efficientWeight = 1
				case codes.ResourceExhausted:
					res.efficientWeight = res.efficientWeight / 2
				case codes.Aborted:
					res.efficientWeight = res.efficientWeight / 2
				default:
					if res.efficientWeight == 1 {
						return
					}
					res.efficientWeight--
				}
			}
		},
	}, nil
}

type weightConn struct {
	balancer.SubConn

	mutex sync.Mutex

	weight          int
	currentWeight   int
	efficientWeight int

	// 可以用来标记不可用（比如说熔断了）
	available bool
}
