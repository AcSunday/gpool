package gpool

import (
	"errors"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

type GoPool interface {
	// FastTune tune the number of goroutines, but not over max size set by SetPoolMaxSize
	FastTune(size int) error
	// SetPoolMaxSize the number of goroutines used for control max size
	SetPoolMaxSize(size int) error
	// SubmitFunc asynchronously submit functions to goroutines for execution
	SubmitFunc(f func()) error
	// GetCurrentGoroutineNum get the current number of running goroutines
	GetCurrentGoroutineNum() int
	// GetCurrentPoolCap get the current number of pool capacity
	GetCurrentPoolCap() int
	// Close pool recycle && release
	Close()
}

type goPool struct {
	minRoutineSize int
	maxRoutineSize int
	criticalValue  float64

	tunePeriod time.Duration
	pool       *ants.Pool
	lock       sync.Locker
}

//
// NewGoPool
//  @param minSize 缩容最小值
//  @param maxSize 扩容最大值
//  @param criticalValue 扩容临界值，缩容是该值的50%
//  @param tuneInterval  tune pool size interval time
//  @param panicHandler  panic handler function
//  @return GoPool
//  @return error
//
func NewGoPool(minSize, maxSize int, criticalValue float64, tuneInterval time.Duration, panicHandler func(interface{})) (GoPool, error) {
	if minSize < 0 || maxSize < minSize {
		return nil, errors.New("invalid size")
	}
	if criticalValue < 0 || criticalValue > 1 {
		return nil, errors.New("invalid critical value")
	}

	pool, err := ants.NewPool(minSize, ants.WithPreAlloc(true), ants.WithPanicHandler(panicHandler))
	if err != nil {
		return nil, err
	}

	g := &goPool{
		minRoutineSize: minSize,
		maxRoutineSize: maxSize,
		tunePeriod:     tuneInterval,
		criticalValue:  criticalValue,
		pool:           pool,
		lock:           &sync.Mutex{},
	}
	go g.autoTune()
	return g, nil
}
