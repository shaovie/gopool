package gopool

import (
	"sync/atomic"
	"time"
)

// GoPool is a minimalistic goroutine pool that provides a pure Go implementation
type GoPool struct {
	noCopy

	queueLen atomic.Int32
	doTaskN  atomic.Int32
	workerN  atomic.Int32
	options  options

	workerSem chan struct{}
	queue     chan func()
}

// NewGoPool provite fixed number of goroutines, reusable. M:N model
//
// M: the number of reusable goroutines,
// N: the capacity for asynchronous task queue.
func NewGoPool(opts ...Option) *GoPool {
	opt := setOptions(opts...)
	if opt.minWorkers <= 0 {
		panic("GoPool: min workers <= 0")
	}
	if opt.minWorkers > opt.maxWorkers {
		panic("GoPool: min workers > max workers")
	}
	p := &GoPool{
		options:   opt,
		workerSem: make(chan struct{}, opt.maxWorkers),
		queue:     make(chan func(), opt.queueCap),
	}
	for i := int32(0); i < p.options.minWorkers; i++ { // pre spawn
		p.workerSem <- struct{}{}
		go p.worker(func() {})
	}
	go p.shrink()
	return p
}

// QueueFree returns (capacity of task-queue - length of task-queue)
func (p *GoPool) QueueFree() int {
	return int(p.options.queueCap - p.queueLen.Load())
}

// Workers returns current the number of workers
func (p *GoPool) Workers() int {
	return int(p.workerN.Load())
}

// Go submits a task to this pool.
func (p *GoPool) Go(task func()) {
	if task == nil {
		panic("GoPool: Go task is nil")
	}
	select {
	case p.queue <- task:
		p.queueLen.Add(1)
	case p.workerSem <- struct{}{}:
		go p.worker(task)
	}
}

func (p *GoPool) worker(task func()) {
	p.workerN.Add(1)
	defer func() {
		<-p.workerSem
		p.workerN.Add(-1)
		if e := recover(); e != nil {
			if p.options.panicHandler != nil {
				p.options.panicHandler(e)
			}
		}
	}()

	for {
		task()
		task = <-p.queue
		if task == nil {
			break
		}
		p.doTaskN.Add(1)
		p.queueLen.Add(-1)
	}
}
func (p *GoPool) shrink() {
	ticker := time.NewTicker(p.options.shrinkPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			doTaskN := p.doTaskN.Load()
			p.doTaskN.Store(0)
			if doTaskN < p.options.tasksBelowN {
				closeN := p.workerN.Load() - p.options.minWorkers
				for closeN > 0 {
					p.queue <- nil
					closeN--
				}
			}
		}
	}
}

// Detecting illegal struct copies using `go vet`
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
