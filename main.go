package main

import (
	"fmt"

	"github.com/apbgo/go-study-group/chapter2"
)

func main() {
	f := chapter2.Fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
