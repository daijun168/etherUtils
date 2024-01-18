package util

import (
	"fmt"
	"testing"
)

func TestShowAddress(t *testing.T) {
	addr := FormatAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")
	fmt.Println(addr)
}
