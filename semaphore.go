package main

import "sync"

type Semaphore struct {
	permits int
	cond    *sync.Cond
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		permits: n,
		cond: sync.NewCond(&sync.Mutex{}),	// new conditon variable & associated mutex on new semaphore
	}
}

func (rw *Semaphore) Acquire() {
	rw.cond.L.Lock()	//acquires mutex to protect permit variables

	for rw.permits <= 0 {
		rw.cond.Wait() // wait until available
	}
	rw.permits--	//dec. available permits
	rw.cond.L.Unlock()	//releases mutex
}

func (rw*Semaphore) Release() {
	rw.cond.L.Lock()
	
	rw.permits ++ // inc. no. of available permit
	rw.cond.Signal()

	rw.cond.L.Unlock()
}