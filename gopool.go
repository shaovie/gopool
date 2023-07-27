package gopool

import (
	"sync/atomic"
)

// GoPool is a minimalistic goroutine pool that provides a pure Go implementation
type GoPool struct {
	noCopy

	queueCap     int
	queueLen     atomic.Int32
	panicHandler func(any)

	workerSem chan struct{}
	queue     chan func()
}

// NewGoPool provite fixed number of goroutines, reusable. M:N model
//
// M: the number of reusable goroutines,
// N: the capacity for asynchronous task queue.
func NewGoPool(sizeM, preSpawn, queueCap int, panicHandler func(any)) *GoPool {
	if preSpawn <= 0 && queueCap > 0 {
		panic("GoPool: dead queue")
	}
	if preSpawn > sizeM {
		preSpawn = sizeM
	}
	p := &GoPool{
		queueCap:     queueCap,
		panicHandler: panicHandler,
		workerSem:    make(chan struct{}, sizeM),
		queue:        make(chan func(), queueCap),
	}
	for i := 0; i < preSpawn; i++ { // pre spawn
		p.workerSem <- struct{}{}
		go p.worker(func() {})
	}
	return p
}

// QueueFree returns (capacity of task-queue - length of task-queue)
func (p *GoPool) QueueFree() int {
	return p.queueCap - int(p.queueLen.Load())
}

// Go submits a task to this pool.
func (p *GoPool) Go(task func()) {
	select {
	case p.queue <- task:
		p.queueLen.Add(1)
	case p.workerSem <- struct{}{}:
		go p.worker(task)
	}
}

func (p *GoPool) worker(task func()) {
	defer func() {
		<-p.workerSem
		if e := recover(); e != nil {
			if p.panicHandler != nil {
				p.panicHandler(e)
			}
		}
	}()

	for {
		task()
		task = <-p.queue
		p.queueLen.Add(-1)
	}
}

// Detecting illegal struct copies using `go vet`
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
