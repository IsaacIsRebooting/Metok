package gofer

import (
	"sync"
	"time"

	"github.com/TremblingV5/box/rearer"
)

var useGlobalPool bool
var setUseGlobalPoolOnce sync.Once

func SetUseGlobalPool(value bool) {
	setUseGlobalPoolOnce.Do(func() {
		useGlobalPool = value
	})
}

func Go(f func()) {
	if useGlobalPool {
		if pool == nil {
			InitGlobalPool()
		}
		_ = pool.Submit(f)
		return
	}
	go func() {
		defer rearer.Recover() //假设 rearer.Recover 捕获并处理 panic
		f()
	}()
}

func GoWithTimeout(f func(), d time.Duration) (isFinish bool) {
	ch := make(chan struct{})

	Go(func() {
		f()
		close(ch)
	})

	timer := time.NewTimer(d)
	select {
	case <-timer.C:
		return false
	case <-ch:
		timer.Stop()
		return true
	}
}
