package util

import (
	"sync"
	"time"
)

type (
	record[T any] struct {
		waiter sync.WaitGroup
		val    T
		err    error
	}
	ProducerFunc[T any] func() (value T, err error)
	Group[T any]        struct {
		mutex   sync.Mutex
		records map[string]*record[T]
	}
	panicErr struct {
		reason any
	}
)

var GoexitCalledError error = &BizErr{
	Msg: "runtime.Goexit() called",
}

var _ error = (*panicErr)(nil)

// Error implements error.
func (p *panicErr) Error() string {
	return "producer cause panic"
}

func (g *Group[T]) Do(key string, producer ProducerFunc[T], withTTL ...time.Duration) (result T, shared bool, err error) {
	ttl := time.Millisecond * 0
	switch {
	case len(withTTL) > 1:
		panic("withTTL should only 1 parameter")
	case len(withTTL) == 1:
		ttl = withTTL[0]
	}
	g.mutex.Lock()

	if g.records == nil {
		g.records = make(map[string]*record[T], 16)
	}
	if exist, ok := g.records[key]; ok {
		g.mutex.Unlock()

		exist.waiter.Wait()
		if exist.err != nil {
			// 将 panic 重新外抛
			if err, ok := exist.err.(*panicErr); ok {
				panic(err.reason)
			}
		}

		return exist.val, true, exist.err
	}

	newRecord := new(record[T])
	newRecord.waiter.Add(1)
	g.records[key] = newRecord
	g.mutex.Unlock()
	g.doCall(key, newRecord, producer, ttl)

	return newRecord.val, false, newRecord.err
}

func (g *Group[T]) doCall(key string, newRecord *record[T], producer ProducerFunc[T], ttl time.Duration) {
	// 表示函数正常返回
	normaReturn := false
	// 表示存在 panic 引起的崩溃
	recovered := false

	// 在最后执行(defer 的执行顺序为 LIFO)
	defer func() {
		// 用户调用了 runtime.Goexit 方法
		if !normaReturn && !recovered {
			// 设置标志
			newRecord.err = GoexitCalledError
		}

		newRecord.waiter.Done()
		// 延迟执行清理
		go func() {
			g.mutex.Lock()
			defer g.mutex.Unlock()
			if ttl > 0 {
				time.Sleep(ttl)
			}
			if g.records[key] == newRecord {
				delete(g.records, key)
			}
		}()

		// 确保初次执行时遇到 panic 就能往外抛, 后续的依靠缓存里的判断
		switch {
		case recovered:
			if r, ok := newRecord.err.(*panicErr); ok {
				panic(r.reason)
			}
		case newRecord.err == GoexitCalledError:
			// 暂不做处理
		}
	}()

	func() {
		// 捕获 producer 可能产生的错误
		defer func() {
			// 没有正常返回, 即存在两种情况: panic, 或者调用了 runtime.Goexit()
			if !normaReturn {
				if panicE := recover(); panicE != nil {
					newRecord.err = &panicErr{
						reason: panicE,
					}
				}

				// runtime.Goexit() 并不会导致 recover 捕获实际的错误(即捕获结果为 nil), 所以依然需要设置恢复标志
				recovered = true
			}
		}()

		// 产生的 error 和 runtime.Goexit() 都被视为正常操作
		newRecord.val, newRecord.err = producer()
		normaReturn = true
	}()
}
