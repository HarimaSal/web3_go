package main

import "fmt"

func main() {
	arr1 := [3]int{1, 2, 3}
	arr2 := [2]int{4, 5}
	merged := append(3, arr2[:]...)
	fmt.Println(merged)
}
