package main

import (
	"errors"
	"fmt"
)

func Slice() {
	s1 := []int{1, 2, 3, 4}
	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))

	// 创建一个长度为 3，容量为 4 的切片
	s2 := make([]int, 3, 4)
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))

	s2 = append(s2, 7)
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))

	// 扩容了
	s2 = append(s2, 8)
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))

	s3 := make([]int, 4)
	fmt.Printf("s3: %v, len: %d, cap: %d \n", s3, len(s3), cap(s3))

	fmt.Printf("s3[2]:%d", s3[2])
	//fmt.Printf("s3[99]:%d", s3[99])
}

func SubSlice() {
	s1 := []int{2, 4, 6, 8, 10}
	s2 := s1[1:3]
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))

	s3 := s1[1:]
	fmt.Printf("s3: %v, len: %d, cap: %d \n", s3, len(s3), cap(s3))

	s4 := s1[:3]
	fmt.Printf("s4: %v, len: %d, cap: %d \n", s4, len(s4), cap(s4))
}

func ShareSlice() {
	s1 := []int{1, 2, 3, 4}
	s2 := s1[2:]
	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))

	s2[0] = 99
	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))

	s2 = append(s2, 199)
	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))
	s2[1] = 1999
	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))
}

func shrink[T any](s []T) ([]T, bool) {
	l, c := len(s), cap(s)
	newCap := 0
	if c <= 64 {
		return s, false
	}
	if c > 2048 && c/l >= 4 {
		factor := 0.625
		newCap = int(float32(c) * float32(factor))
	}
	if c <= 2048 && c/l >= 2 {
		newCap = c / 2
	}

	s1 := make([]T, l, newCap)
	s1 = append(s1, s...)
	return s1, true
}

func Delete[T any](idx int, vals []T) ([]T, error) {

	if idx < 0 || idx >= len(vals) {
		return []T{}, errors.New("index out of range")
	}

	ret := append(vals[:idx], vals[idx+1:]...)
	ret, _ = shrink(ret)

	return ret, nil
}
