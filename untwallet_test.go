package main

import (
	"fmt"
	"testing"
)

func TestUntWallet(t *testing.T) {
	wallet := UntWallet{}
	wallet.GetToken()
	fmt.Println(wallet.token)
}
