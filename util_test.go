package util

import (
	"fmt"
	"testing"
)

func TestShowAddress(t *testing.T) {
	addr := FormatAddress("0xfffffffffffffffffff")
	fmt.Println(addr)

	address, k, err := GetAddAndKey("----")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(address, k)
}
