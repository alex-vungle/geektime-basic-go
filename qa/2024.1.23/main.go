package main

import (
	"fmt"
)

func main() {
	Defer()
}

func Defer() {
	i := 0
	fmt.Printf("print(i): %v\n", i)
	defer fmt.Printf("defer 直接调用 print(i): %v\n", i)

	//defer func(i int) {
	//	fmt.Printf("defer 直接形参调用 print(i): %v\n", i)
	//}(i)

	defer func() {
		fmt.Printf("defer 闭包函数 print(i): %v\n", i)
	}()

	i++
	fmt.Printf("print(i): %v\n", i)

	// runtime.NumGoroutine()
	// 如果你接入了 prometheus，它默认就会采集 goroutine 数量

	// var db *sql.DB
	// db.Stats()
}

//var _ error = myError{}

//type myError struct {
//}

//func (m myError) Error() string {
//
//}

// 假如说我在这里用了额
//func Mybiz() error{
//	return &myError{}
//}
