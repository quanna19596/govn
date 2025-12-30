package main

import "fmt"

func main() {
	a := 5
	b := &a
	c := *b

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
}
