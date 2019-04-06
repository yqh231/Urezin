package main

import (
	"time"

	"github.com/yqh231/Urezin/pool"
)

var Pool *pool.GoPool

func main() {
	untWallet := NewUntWallet()
	Pool = pool.NewPool(1, 10 * time.Second)
	go Pool.Run(untWallet)
	Pool.RunTillRoutineShut()
}
