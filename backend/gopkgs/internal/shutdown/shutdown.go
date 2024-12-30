package shutdown

import (
	"os"
	"os/signal"
	"sort"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/IsaacIsRebooting/Metok/backend/gopkgs/gofer"
	"github.com/samber/lo"
)

var globalShutdownManager = newManager()

// init()函数，Go程序在每个包加载时默认调用，可以用来执行包级别的初始化任务。
func init() {
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
		sig := <-signals
		globalShutdownManager.Fire(sig)
		signal.Stop(signals)
	}()
}

type orderHandler struct {
	fns []func()
}
type manager struct {
	mu         sync.RWMutex
	fired      atomic.Bool
	firedCh    chan struct{}
	doneCh     chan struct{}
	lastSignal atomic.Value
	orderMap   map[int]*orderHandler
}

func newManager() *manager {
	return &manager{
		firedCh:  make(chan struct{}),
		doneCh:   make(chan struct{}),
		orderMap: make(map[int]*orderHandler),
	}
}

func (m *manager) Fire(sig os.Signal) {
	if !m.fired.CompareAndSwap(false, true) {
		return
	}
	close(m.firedCh)
	m.lastSignal.Store(sig)

	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := lo.Keys(m.orderMap)
	sort.Ints(keys)

	for _, key := range keys {
		m.fireOrder(m.orderMap[key])
	}
	close(m.doneCh)
}

func (m *manager) fireOrder(order *orderHandler) {
	wg := sync.WaitGroup{}

	for _, fn := range order.fns {
		wg.Add(1)
		tmp := fn
		gofer.Go(func() {
			defer wg.Done()
			tmp()
		})
		wg.Wait()
	}
}

func (m *manager) AddOrderHandler(order int, fn func()) {
	if m.fired.Load() {
		gofer.Go(fn)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	orderHandle, ok := m.orderMap[order]
	if !ok {
		orderHandle = &orderHandler{}
		m.orderMap[order] = orderHandle
	}

	orderHandle.fns = append(orderHandle.fns, fn)
}

func (m *manager) LastSignal() os.Signal {
	r, ok := m.lastSignal.Load().(os.Signal)
	if !ok {
		return nil
	}
	return r
}

func (m *manager) FiredCh() <-chan struct{} {
	return m.firedCh
}

func (m *manager) Wait(duration time.Duration) {
	if duration <= 0 {
		<-m.firedCh
		return
	}
	_ = gofer.GoWithTimeout(func() {
		<-m.firedCh
	}, duration)
}



func Wait(duration time.Duration) {
	globalShutdownManager.Wait(duration)
}

func FiredCh() <-chan struct{} {
	return globalShutdownManager.FiredCh()
}
