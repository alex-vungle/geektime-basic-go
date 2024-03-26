package main_test

import (
	"fmt"
	"testing"
)

func TestDeleteQA(t *testing.T) {
	var s = []int{1, 2, 4, 5, 3} // 使用变量替换[]int{1, 2, 4, 5, 3}有问题 原因？？

	// 1.
	fmt.Printf("%v, %v \n", Delete[int](2, s), cap(Delete[int](2, s))) //[1 1 5 3], 4

	fmt.Printf("%v, %v \n", Delete[int](0, s), cap(Delete[int](0, s))) //[1 1 5 3], 4
	fmt.Printf("%v, %v \n", Delete[int](1, s), cap(Delete[int](1, s))) //[1 1 5 3], 4
	fmt.Println("######")

	fmt.Printf("%v, %v \n", Delete[int](2, []int{1, 2, 4, 5, 3}), cap(Delete[int](2, []int{1, 2, 4, 5, 3})))
	fmt.Printf("%v, %v \n", Delete[int](4, []int{1, 2, 4, 5, 3}), cap(Delete[int](4, []int{1, 2, 4, 5, 3})))
	fmt.Printf("%v, %v \n", Delete[int](0, []int{1, 2, 4, 5, 3}), cap(Delete[int](0, []int{1, 2, 4, 5, 3})))
	fmt.Printf("%v, %v \n", Delete[int](1, []int{1, 2, 4, 5, 3}), cap(Delete[int](1, []int{1, 2, 4, 5, 3})))
	fmt.Printf("%v, %v \n", Delete[int](3, []int{1, 2, 4, 5, 3}), cap(Delete[int](3, []int{1, 2, 4, 5, 3})))

	var sany any = s
	fmt.Printf("%v", sany)

}

func Delete[T any](index int, slice []T) []T {

	if index < 0 || index > len(slice)-1 {
		fmt.Println("index error")
	}
	for i := index; i > 0; i-- {
		slice[i] = slice[i-1]
	}
	slice = slice[1:]

	return slice
}
