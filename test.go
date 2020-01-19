package main

import "fmt"

func main() {
	a := make(map[int32] bool, 0)
	a[10] = true
	fmt.Println(a[1])
	if _, ok := a[1]; ok != true{
		fmt.Println("sdasdsadsa")
	}
}
