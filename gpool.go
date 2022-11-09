package gpool

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"time"
)

// 扩容数量
func (g *goPool) calcExpansionNum() int {
	var ret = 10
	return ret
}

// 缩减数量
func (g *goPool) calcReduceNum() int {
	var ret = 10
	return ret
}

// 自动扩容
func (g *goPool) autoTune() {
	var ticker = time.NewTicker(g.tunePeriod)
	defer ticker.Stop()

	for range ticker.C {
		free := g.pool.Running()
		poolCap := g.pool.Cap()
		percentUsed := float64(free) / float64(g.maxRoutineSize)
		reduceCriticalValue := g.criticalValue / 2

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
	return g.pool.Free()
}

func (g *goPool) Close() {
	g.pool.Release()
}
