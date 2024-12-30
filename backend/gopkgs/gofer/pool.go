package gofer

import (
	"github.com/TremblingV5/box/rearer"
	"github.com/panjf2000/ants/v2"
)

var (
	pool         *Pool
	panicHandler = ants.WithPanicHandler(func(err any) {
		rearer.LogRecoverStack(err)
	})
)

var defaultPoolSize = 1000

func SetPoolSize(size int) {
	defaultPoolSize = size
}

type Pool struct {
	poolWithFunc ants.PoolWithFunc
	pool         *ants.Pool
}

func InitGlobalPool() {
	p, _ := ants.NewPool(defaultPoolSize, panicHandler)
	pool = &Pool{
		pool: p,
	}
}

// Submit 提交一个任务到任务池中执行。
// 该方法首先检查任务池是否已经初始化，如果未初始化，则进行初始化。
// 然后，尝试将给定的任务提交到任务池中执行。
// 如果任务提交成功，返回nil；如果任务提交失败，返回错误。
func (p *Pool) Submit(task func()) error {
    // 检查任务池是否已经初始化，如果未初始化，则进行初始化。
    if p.pool == nil {
        pool, _ := ants.NewPool(defaultPoolSize, panicHandler)
        p.pool = pool
    }
    // 尝试将给定的任务提交到任务池中执行。
    if err := p.pool.Submit(task); err != nil {
        return err
    }
    // 任务提交成功，返回nil。
    return nil
}
