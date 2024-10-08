package main

import (
	"container/list"
	"sync"
)

type Channel[M any] struct {
	capacitySema Semaphore	//blocks sender when buffer is full
	sizeSema     Semaphore	//buff. size semaphore to block receiver when buffer is empty	
	mutex        sync.Mutex	//protecting shared resources
	buffer       *list.List
}

func NewChannel [M any] (capacity int) *Channel[M] {
	return &Channel[M]{
		// new semaphore with no. of permits equal to input capacity
		capacitySema: *NewSemaphore(capacity),
		
		//new semaphore with no. of permits equal to 0 
		sizeSema: *NewSemaphore(0),
		buffer: list.New(),	// new empty linked list
	}
}

func (c *Channel[M]) Send(message M) {
	c.capacitySema.Acquire()

	c.mutex.Lock()
	c.buffer.PushBack(message)
	c.mutex.Unlock()

	c.sizeSema.Release()
}

func (c *Channel[M]) Receive() M {
	c.capacitySema.Release()	//release one permit from capacity semaphore
	c.sizeSema.Acquire()	// acquire one permit from buffer size semaphore
	
	c.mutex.Lock()
	v := c.buffer.Remove(c.buffer.Front()).(M)	//removes one message from buffer while using mutex to protect against race conditions
	c.mutex.Unlock()

	return v
}