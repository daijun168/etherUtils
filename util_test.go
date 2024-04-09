package util

import (
	"fmt"
	"testing"
)

func TestShowAddress(t *testing.T) {
	addr := FormatAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")
	fmt.Println(addr)

	address, k, err := GetAddAndKey("0xBc1439D5bFCCcA724C5C864916ce55eC248cbc7C----b38986bc3f330a3e65c327016a6a8fb27e701c6505ffb2de2549d7b10a626075")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(address, k)
}
