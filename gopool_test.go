package gopool

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

var cnt atomic.Int32

func catchPanic(p any) {
	//fmt.Printf("catch panic: %v\n", p)
}
func job() {
	cnt.Add(1)
	defer cnt.Add(-1)
	v := rand.Int63() % 10
	if v == 0 {
		panic("panic")
	}
	time.Sleep(time.Duration(v) * time.Millisecond)
}
func TestGoPool(t *testing.T) {
	fmt.Println("hello boy")
	p := NewGoPool(
		MinWorkers(16),
		MaxWorkers(1024),
		QueueCap(2048),
		ShrinkPeriod(time.Millisecond*1500),
		TasksBelowNToShrink(12800),
		PanicHandler(catchPanic),
	)
	time.Sleep(100 * time.Millisecond)
	fmt.Println("queue free:", p.QueueFree(), "workers:", p.Workers())
	for i := 0; i < 2049; i++ {
		if p.QueueFree() > 0 {
			p.Go(job)
		}
	}

	fmt.Println("queue free:", p.QueueFree(), "workers:", p.Workers())
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("cnt:", cnt.Load(), "queue free:", p.QueueFree(), "workers:", p.Workers())
	time.Sleep(600 * time.Millisecond)
	fmt.Println("cnt:", cnt.Load(), "queue free:", p.QueueFree(), "workers:", p.Workers())
}
