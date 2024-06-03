package main

import "fmt"

func main() {

	a := []int{1, 2, 3, 4}

	fmt.Println("Before: ", cap(a))

	result, err := Delete(0, a)

	fmt.Println("After the 1st operation:", a)

	result, err = Delete(0, result)

	fmt.Println("After the 2nd operation:", a)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}

	fmt.Println("After:", cap(a))
}
