package pools

import (
	"fmt"
	"sync"
)

// WorkloadLock 线程总量互斥锁
type WorkloadLock struct {
	Current int         // 线程总量进度
	Locker  *sync.Mutex // 互斥锁
}

// Pool 线程池设置
type Pool struct {
	Routines int    // 线程数量
	WorkLoad int    // 处理数量
	Callback func() // 处理事件
}

var wg sync.WaitGroup

var chs chan int

// Work 进程并发处理
func (p *Pool) Work(wg *sync.WaitGroup) {

	chs = make(chan int, 10)
	WorkloadLock := new(WorkloadLock)

	for i := 1; i <= p.Routines; i++ {
		wg.Add(1)
		go p.run(WorkloadLock, wg, chs)
	}

	go func() {
		wg.Wait()
		close(chs)
	}()

	for i := 1; i <= p.Routines; i++ {
		<-chs
	}
}

func (p *Pool) run(work *WorkloadLock, wg *sync.WaitGroup, ch chan int) {

	work.Locker = new(sync.Mutex)

	for work.Current < p.WorkLoad {

		work.Locker.Lock()
		work.Current++
		private := work.Current
		work.Locker.Unlock()

		fmt.Println(private)
	}

	ch <- 1
	wg.Done()
}
