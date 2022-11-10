package gpool

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"time"
)

// 扩容数量
func (g *goPool) calcExpansionNum() int {
	var newcap = g.pool.Cap()
	doublecap := newcap + newcap

	if g.pool.Cap() < 1024 {
		newcap = doublecap
	} else {
		newcap += newcap / 4

		// 超过最大限制 或是 溢出
		if newcap > g.maxRoutineSize || newcap <= 0 {
			newcap = g.maxRoutineSize
		}
	}

	return newcap
}

// 缩减数量
func (g *goPool) calcReduceNum() int {
	var newcap = g.pool.Cap() / 4

	// 小于最小限制
	if newcap < g.minRoutineSize {
		newcap = g.minRoutineSize
	}

	return newcap
}

// 自动扩容
func (g *goPool) autoTune() {
	var ticker = time.NewTicker(g.tunePeriod)
	defer ticker.Stop()

	for range ticker.C {
		runNum := g.pool.Running()
		poolCap := g.pool.Cap()
		percentUsed := float64(runNum) / float64(poolCap)
		reduceCriticalValue := g.criticalValue * 0.5

		// 占比大于临界值，并且当前池大小 < 可扩容的最大值，则扩容池大小
		if percentUsed > g.criticalValue && poolCap < g.maxRoutineSize {
			if err := g.FastTune(g.calcExpansionNum()); err != nil {
				log.WithError(err).Errorf("[gpool] expansion tune size fail")
			}
		} else if percentUsed < reduceCriticalValue && poolCap > g.minRoutineSize {
			// 占比小于缩减临界值，并且当前池大小 > 池缩容的最小值，则缩减池大小
			if err := g.FastTune(g.calcReduceNum()); err != nil {
				log.WithError(err).Errorf("[gpool] reduce tune size fail")
			}
		} else if poolCap > g.maxRoutineSize {
			// 当前池大小 超过 可扩容的最大值，则缩减到最大值
			if err := g.FastTune(g.maxRoutineSize); err != nil {
				log.WithError(err).Errorf("[gpool] the pool size more than max size, reduce tune size fail")
			}
		}

	}
}

func (g *goPool) FastTune(size int) error {
	if size > g.maxRoutineSize {
		return errors.New("over max")
	} else if size < g.minRoutineSize {
		return errors.New("over min")
	}

	g.lock.Lock()
	defer g.lock.Unlock()
	g.pool.Tune(size)
	return nil
}

func (g *goPool) SetPoolMaxSize(size int) error {
	if size < g.minRoutineSize {
		return errors.New("must be greater than the minimum size")
	}

	g.lock.Lock()
	defer g.lock.Unlock()
	g.maxRoutineSize = size
	return nil
}

func (g *goPool) SubmitFunc(f func()) error {
	return g.pool.Submit(f)
}

func (g *goPool) GetCurrentGoroutineNum() int {
	return g.pool.Running()
}

func (g *goPool) GetCurrentPoolCap() int {
	return g.pool.Cap()
}

func (g *goPool) Close() {
	g.pool.Release()
}
