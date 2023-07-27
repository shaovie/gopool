package gopool

import (
	"fmt"
    "time"
    "math/rand"
	"sync/atomic"
	"testing"
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
    p := NewGoPool(512, 128, 1024, catchPanic)
    time.Sleep(100 * time.Millisecond)
    fmt.Println("queue free:", p.QueueFree())
    for i := 0; i < 2048; i++ {
        p.Go(job)
    }

    fmt.Println("queue free:", p.QueueFree())
    time.Sleep(1000 * time.Millisecond)
    fmt.Println("cnt:", cnt.Load())
}
