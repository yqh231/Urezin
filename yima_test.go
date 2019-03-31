package main

import (
	"testing"

	"fmt"
)

func TestYiMa(t *testing.T) {
	ym := NewYiMa()

	ym.Login("yqh231", "2316678")
	fmt.Println(ym.token)
	phone := ym.GetPhone("21714", "")
	fmt.Println(phone)
}
