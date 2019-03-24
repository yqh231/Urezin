package pool

import (
	"sync"
	"sync/atomic"
	"time"
	"fmt"

	"github.com/yqh231/Urezin/log"
)

type Worker interface {
	Work()
}

type GoPool struct {
	wg sync.WaitGroup
	kill chan struct{}
	shutdown chan struct{}
	addRoutine chan struct{}
	removeRoutine chan struct{}
	tasks chan Worker
	active int64
	routines int64
	pending int64
	timeDisplay time.Duration
}


func NewPool(routines int64, duration time.Duration) *GoPool{
	pool := &GoPool{
		kill: make(chan struct{}),
		shutdown: make(chan struct{}),
		addRoutine: make(chan struct{}),
		removeRoutine: make(chan struct{}),
		tasks: make(chan Worker),
		timeDisplay: duration,
	}
	pool.manager()
	pool.Add(int(routines))
	return pool
}

func(gp *GoPool) work(){
done:
	for {
		select {
		case t := <- gp.tasks:
			atomic.AddInt64(&gp.active, 1)
			t.Work()
			atomic.AddInt64(&gp.active, -1)
		case <- gp.kill:
			log.Info.Println("shut down routine")
			break done
		}
	}
	atomic.AddInt64(&gp.pending, -1)
	gp.wg.Done()
}

func(gp *GoPool) Shutdown() {
	close(gp.shutdown)
	gp.wg.Wait()
}

func(gp *GoPool) RunTillRoutineShut(){
	gp.wg.Wait()
}

func(gp *GoPool) Run(work Worker){
	atomic.AddInt64(&gp.pending, 1)
	gp.tasks <- work
	atomic.AddInt64(&gp.pending, -1)
}

func(gp *GoPool) Add(routines int){
	for i := 0; i < routines; i++{
		gp.addRoutine <- struct{}{}
	}
}

func(gp *GoPool) manager(){

	gp.wg.Add(1)

	go func(){
	
	timer := time.NewTimer(gp.timeDisplay)

	for {
		select {
		case <- gp.shutdown:
			routines := int(atomic.LoadInt64(&gp.routines))

			for i := 0; i < routines; i++{
				gp.kill <- struct{}{}
			}
			gp.wg.Done()
			return
		
		case <- gp.addRoutine:
			gp.wg.Add(1)
			atomic.AddInt64(&gp.routines, 1)
			go gp.work()
		case <- gp.removeRoutine:	
			atomic.AddInt64(&gp.routines, -1)
			gp.kill <- struct{}{}
		case <- timer.C:
			routines := atomic.LoadInt64(&gp.routines)
			log.Info.Println(fmt.Sprintf(
				"pool status, total routines %v ,active routines %v, pending routines %v", routines, gp.active, gp.pending))
			timer.Reset(gp.timeDisplay)
		}
	}
	}()

}