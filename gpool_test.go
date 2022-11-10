package gpool

import (
	"testing"
	"time"
)

// go test -bench . -benchmem -benchtime=5s -count=3 -timeout=30s

const calcTimes = 30

var g GoPool

func init() {
	g, _ = NewGoPool(30, 2048, 0.5, time.Millisecond*100, nil)
}

func fib(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return fib(n-2) + fib(n-1)
}

func TestCalcFibSpendTime(t *testing.T) {
	start := time.Now()
	fib(calcTimes)
	spend := time.Now().Sub(start)
	t.Logf("spend %v", spend)
}

func BenchmarkGoPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g.SubmitFunc(func() {
			fib(calcTimes)
		})
	}
}

func BenchmarkGoroutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go fib(calcTimes)
	}
}
