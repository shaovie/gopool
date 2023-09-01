package gopool

import (
	"time"
)

// options provides all optional parameters
type options struct {
	queueCap     int32
	minWorkers   int32
	maxWorkers   int32
	tasksBelowN  int32 // in shrinkPreiod
	shrinkPeriod time.Duration
	panicHandler func(any)
}

// Option function
type Option func(*options)

func setOptions(optL ...Option) options {
	opts := options{
		queueCap:     128,
		minWorkers:   8,
		maxWorkers:   256,
		tasksBelowN:  1024,
		shrinkPeriod: time.Minute,
	}

	for _, opt := range optL {
		opt(&opts)
	}
	return opts
}

// QueueCap set the capacity of the pool's queue
func QueueCap(v int32) Option {
	return func(o *options) {
		if v < 1 {
			panic("gopool:QueueCap: param is illegal")
		}
		o.queueCap = v
	}
}

// MinWorkers set min workers
func MinWorkers(v int32) Option {
	return func(o *options) {
		if v < 1 {
			panic("gopool:MinWorkers: param is illegal")
		}
		o.minWorkers = v
	}
}

// MaxWorkers set max workers
func MaxWorkers(v int32) Option {
	return func(o *options) {
		if v < 1 {
			panic("gopool:MinWorkers: param is illegal")
		}
		o.maxWorkers = v
	}
}

// ShrinkPeriod set shrink cycle
func ShrinkPeriod(v time.Duration) Option {
	return func(o *options) {
		if v < 1 {
			panic("gopool:ShrinkPeriod: param is illegal")
		}
		o.shrinkPeriod = v
	}
}

// TasksBelowNToShrink set shrink condition
func TasksBelowNToShrink(v int32) Option {
	return func(o *options) {
		if v < 1 {
			panic("gopool:TasksBelowNToShrink: param is illegal")
		}
		o.tasksBelowN = v
	}
}

// PanicHandler set panic handler
func PanicHandler(fn func(any)) Option {
	return func(o *options) {
		if fn == nil {
			panic("gopool:PanicHandler: param is illegal")
		}
		o.panicHandler = fn
	}
}
