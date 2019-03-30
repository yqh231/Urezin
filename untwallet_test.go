package main

import (
	"testing"
	"fmt"
)

func TestUntWallet(t *testing.T){
	wallet := UntWallet{}
	wallet.GetToken()
	fmt.Println(wallet.token)
}