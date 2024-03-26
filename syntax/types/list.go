package main

type List interface {
	Add(index int, val any)
	Append(val any) error
	Delete(index int) error
}

type LinkedList struct {
	head node
	//Head node
}

func (l *LinkedList) Append(val any) error {
	//TODO implement me
	panic("implement me")
}

func (l *LinkedList) Delete(index int) error {
	//TODO implement me
	panic("implement me")
}

func (l *LinkedList) Add(index int, val any) {
	// 实现这个方法
}

type node struct {
	//next node
	next *node
}

func UseListV1() {
	l := &LinkedList{}
	l.Add(1, 123)
	l.Add(1, "123")
	l.Add(1, nil)
}

func IsList(val any) bool {
	// 类型断言
	list, ok := val.(List)
	if ok {
		// 说明 val 真是的是 List
		list.Append(123)
	}

	// 如果 val 不是 List，这里会 panic
	list = val.(List)

	return ok
}

func IsListV1[T any](val T) bool {
	// 没有办法对类型参数断言
	//list, ok := val.(List)
	//if ok {
	// 说明 val 真是的是 List
	//	list.Append(123)
	//}

	// 如果 val 不是 List，这里会 panic
	//list = val.(List)

	// 纯纯 GO 设计导致的垃圾代码
	var anyVal any = val

	// 编译器只支持对any 进行断言
	list, ok := anyVal.(List)
	if ok {
		list.Append(123)
	}

	return ok
}
