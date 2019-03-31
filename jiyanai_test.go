package main

import (
	"fmt"
	"testing"
)

func TestJiYanAi(t *testing.T) {
	jiyan := NewJiYan()

	pic := jiyan.GetPicture("http://untwallet.com/rucaptcha/")
	fmt.Println(jiyan.Distinguish(pic, "2000"))
}
