package main

import (
	"fmt"
	"testing"

	"github.com/wlcy/tron/explorer/core/utils"
)

func TestA(*testing.T) {

	a := make(chan struct{}, 10)
	fmt.Println(len(a))
	for i := 0; i < 10; i++ {
		a <- struct{}{}
	}

	fmt.Println(len(a))
	<-a
	fmt.Println(len(a))
	<-a
	fmt.Println(len(a))

	a <- struct{}{}

	fmt.Println(len(a))

	b(3)

	fmt.Println(utils.HexEncode([]byte{}))
	fmt.Println(utils.ConverTimestamp(1536348924123))
}

func b(i int) bool {
	fmt.Printf("b %v run\n", i)
	defer c(i)

	if i > 0 {
		return b(i - 1)
	}
	return true
}

func c(i int) {
	fmt.Printf("c %v call\n", i)
}
