package main

import "fmt"

func main() {
	a := []int{1, 2, 3, 4, 5, 6, 7}
	//val, err := DeleteSliceByIndex(a, 3)
	val, err := DeleteSlice(a, 3)
	fmt.Println("原切片的容量", cap(a))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("新切片", val)
		fmt.Println("新切片的容量", cap(val))
	}
}
