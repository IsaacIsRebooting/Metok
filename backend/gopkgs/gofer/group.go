package gofer

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
)

// Group 是一个用于管理和同步多个goroutine的结构体。
// 它提供了启动多个goroutine、等待它们全部完成以及处理错误的能力。
type Group struct {
	// wg 是 sync.WaitGroup 的别名，用于等待所有goroutine完成。
	wg sync.WaitGroup

	// isErrorGroup 表示这个 Group 是否应该在出现第一个错误时停止所有goroutine。
	isErrorGroup bool

	// errOnce 确保在所有goroutine中只报告一次错误。
	errOnce sync.Once
	// err 保存从任何goroutine中报告的第一个错误。
	err error

	// numG 是预计在 Group 中运行的goroutine数量。
	numG int
	// gCount 记录已经完成的goroutine数量。
	gCount int64

	// queueSize 表示任务队列的大小，即同时可以运行的任务数量。
	queueSize int
	// queue 是一个用于排队任务的通道，确保并发执行的goroutine数量不超过queueSize。
	queue chan func() error

	// cancel 是一个函数，用于取消 Group 中所有正在运行的goroutine。
	cancel func()
	// ctx 是 Group 的上下文，用于传递取消信号。
	ctx context.Context
	// waitFired 用于标记是否已经调用了 Wait 方法。
	waitFired atomic.Bool
}

// NewGroup 创建一个新的Group实例，并根据提供的选项进行配置。
// 该函数接受一个context.Context参数，用于控制Group的生命周期。
// GroupOption是用于配置Group的函数类型参数，可以有零个或多个。
// 返回值是一个指向Group的指针，表示新创建的Group实例。
func NewGroup(ctx context.Context, options ...GroupOption) *Group {
	// 初始化Group实例，包含一个sync.WaitGroup用于等待所有任务完成。
	g := &Group{
		wg: sync.WaitGroup{},
	}

	// 遍历所有提供的配置选项，并应用到Group实例上。
	for _, option := range options {
		option(g)
	}

	// 如果Group配置为ErrorGroup，则进一步进行配置。
	if g.isErrorGroup {
		// 如果未设置并发数，则默认为系统的CPU核心数。
		if g.numG == 0 {
			g.numG = runtime.NumCPU()
		}

		// 如果未设置队列大小，则默认为并发数。
		if g.queueSize == 0 {
			g.queueSize = g.numG
		}

		// 创建一个可取消的context，用于后续任务的控制。
		ctx, cancel := context.WithCancel(ctx)
		g.ctx = ctx
		g.cancel = cancel

		// 创建一个错误队列，用于收集并发执行任务时产生的错误。
		g.queue = make(chan func() error, g.queueSize)
	}

	// 返回配置好的Group实例。
	return g
}

// Wait 等待组内的所有协程完成。
// 如果是错误组，先尝试关闭队列，然后等待所有协程完成，最后取消错误组。
func (g *Group) Wait() error {
	// 如果是错误组，并且waitFired标志还未设置，则关闭队列。
	if g.isErrorGroup {
		if g.waitFired.CompareAndSwap(false, true) {
			close(g.queue)
		}
	}

	// 等待所有协程完成。
	g.wg.Wait()

	// 如果是错误组，取消错误组。
	if g.isErrorGroup {
		g.cancel()
	}

	return nil
}
