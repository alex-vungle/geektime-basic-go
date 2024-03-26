package main

import "fmt"

func main() {
	//DeferClosureLoopV1()
	//DeferClosureLoopV2()
	DeferClosureLoopV3()
	//OfNullable[User](User{}).Apply(func(t User) {
	//	println(t.Name)
	//})
}

func DeferClosureLoopV1() {
	for i := 0; i < 10; i++ {
		fmt.Printf("循环 %p \n", &i)
		defer func() {
			// 很多人会预期打出 9 8 7 6 ... 1
			// 最终都是 10 ... 10
			fmt.Printf("%p \n", &i)
			println(i)
		}()
	}
	println("跳出循环")
}

func DeferClosureLoopV2() {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			fmt.Printf("%p \n", &val)
			println(val)
		}(i)
	}
	println("跳出循环")
}

func DeferClosureLoopV3() {
	for i := 0; i < 10; i++ {
		// j 是局部变量。每次循环新建一个
		j := i
		fmt.Printf("循环 %p, %p \n", &i, &j)
		defer func() {
			fmt.Printf("%p \n", &j)
			println(j)
		}()
	}
	println("跳出循环")
}

// 0,0, 1, 1, 2, 3, 4, 10, 10, 10
// 0, 1, 0, 2, 0, 3
// 10, 9, 8, 7
func DeferClosureLoopV4() {
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Printf("%p \n", &i)
			println(i)
		}()
	}
}
