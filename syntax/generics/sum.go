package main

import (
	"io"
)

// T 就是类型参数，T 被约束到了“必须是” Number
// 1.
func Sum[T Number](vals []T) T {
	var res T
	for _, v := range vals {
		res = res + v
	}
	return res
}

func Max[T Number](vals []T) T {
	t := vals[0]
	for i := 1; i < len(vals); i++ {
		if t < vals[i] {
			t = vals[i]
		}
	}
	return t
}

func Min[T Number](vals []T) T {
	t := vals[0]
	for i := 1; i < len(vals); i++ {
		if t > vals[i] {
			t = vals[i]
		}
	}
	return t
}

func Find[T any](vals []T, filter func(t T) bool) T {
	for _, v := range vals {
		if filter(v) {
			return v
		}
	}
	var t T
	return t
}

func Insert[T any](idx int, val T, vals []T) []T {
	if idx < 0 || idx > len(vals) {
		panic("idx不合法")
	}

	// 先扩容
	vals = append(vals, val)
	// 这个写法
	for i := len(vals) - 1; i > idx; i-- {
		if i-1 >= 0 {
			vals[i] = vals[i-1]
		}
	}
	vals[idx] = val
	return vals
}

// Integer 是 int 的衍生类型
type Integer int

// Number 是一个泛型约束
// int 的衍生类型, uint, int32
type Number interface {
	~int | uint | int32
}

func UseSum() {
	res := Sum[int]([]int{123, 123})
	println(res)
	resV1 := Sum[Integer]([]Integer{123, 123})
	println(resV1)
}

// T 被约束为 io.Closer
// io.Closer 是一个普通接口，所以意味着 T 必须实现了 io.Closer 这个接口
func Closable[T io.Closer](t T) {
	t.Close()
}

func TestMyResource() {
	Closable[*myResource](&myResource{})
	// 这种就是类型推断
	// 注意，不是所有情况下都能推断出来
	// 如果你发现编译报错了，就不要依赖于类型推断
	Closable(&myResource{})
}

type myResource struct {
}

func (m *myResource) Close() error {

	//TODO implement me
	panic("implement me")
}
