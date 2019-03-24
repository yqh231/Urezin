package pool

import (
	"testing"
	"fmt"
	"time"
)

type task struct {

}

func(t *task) Work(){
	fmt.Println("test running")
	time.Sleep(time.Second * 3)
}

func TestPool(t *testing.T){
	work := &task{}
	pool := NewPool(5, time.Second * 2)
	for i:=0; i< 100; i++{
		go pool.Run(work)

	}
	pool.Shutdown()
	pool.RunTillRoutineShut()
}