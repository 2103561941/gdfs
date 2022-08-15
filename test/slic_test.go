package main

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	arr := make([]int, 1)
	arr[0] = 1
	s := arr[1:]
	fmt.Println(s, len(s), arr[1:])
}
