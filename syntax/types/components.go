package main

import "io"

type NameI interface {
	Name() string
}

type Outer struct {
	Inner
}

type Outer1 struct {
	*Inner
}

// 如果我调用 Hello 方法，我打印出来的是什么？
// 1. 如果是多态，打印出来的就是 Hello, outer
// 2. 组合的情况下，打印出来的还是 Hello, inner
func (i Outer1) Name() string {
	return "outer"
}

// 组合下
//func (i Outer1) Hello() {
//	println("hello, 我是", i.Inner.Name())
//}

// 这么写就能打出 hello, outer
//func (i Outer1) Hello() {
//	println("hello, 我是", i.Name())
//}

type Inner struct {
}

func (i Inner) Name() string {
	return "inner"
}

func (i Inner) Hello() {
	println("hello, 我是", i.Name())
}

type Outer2 struct {
	io.Closer
}

func Components() {
	var o Outer
	o.Hello()

	//var o1 Outer1
	o1 := Outer1{
		Inner: &Inner{},
	}
	o1.Hello()
}
