package gopool

import (
	"time"
)

// Options provides all optional parameters
type Options struct {
	queueCap     int32
	minWorkers   int32
	maxWorkers   int32
	tasksBelowN  int32 // in shrinkPreiod
	shrinkPeriod time.Duration
	panicHandler func(any)
}

// Option function
type Option func(*Options)

func setOptions(optL ...Option) Options {
	opts := Options{
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
	return func(o *Options) {
		if v > 0 {
			o.queueCap = v
		}
	}
}

// MinWorkers set min workers
func MinWorkers(v int32) Option {
	return func(o *Options) {
		if v > 0 {
			o.minWorkers = v
		}
	}
}

// MaxWorkers set max workers
func MaxWorkers(v int32) Option {
	return func(o *Options) {
		if v > 0 {
			o.maxWorkers = v
		}
	}
}

// ShrinkPeriod set shrink cycle
func ShrinkPeriod(v time.Duration) Option {
	return func(o *Options) {
		if v > 0 {
			o.shrinkPeriod = v
		}
	}
}

// TasksBelowNToShrink set shrink condition
func TasksBelowNToShrink(v int32) Option {
	return func(o *Options) {
		if v > 0 {
			o.tasksBelowN = v
		}
	}
}

// PanicHandler set panic handler
func PanicHandler(fn func(any)) Option {
	return func(o *Options) {
		if fn != nil {
			o.panicHandler = fn
		}
	}
}
