package main

import "fmt"

func main() {
	var nilSlice []int
	emptySlice := []int{}

	fmt.Println("nil slice:")
	fmt.Printf("slice: %v, len: %d, cap: %d\n", nilSlice, len(nilSlice), cap(nilSlice))
	if nilSlice == nil {
		fmt.Println("nilSlice is nil")
	}

	fmt.Println("\nempty slice:")
	fmt.Printf("slice: %v, len: %d, cap: %d\n", emptySlice, len(emptySlice), cap(emptySlice))
	if emptySlice == nil {
		fmt.Println("emptySlice is nil")
	}

	fmt.Println("\nappending to slices:")
	nilSlice = append(nilSlice, 1)
	emptySlice = append(emptySlice, 1)
	fmt.Printf("nilSlice after append: %v\n", nilSlice)
	fmt.Printf("emptySlice after append: %v\n", emptySlice)
}
