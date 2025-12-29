package aiMng

import (
	"sync"
	"sync/atomic"
	"time"
)

// 全局互斥锁和执行状态控制
var (
	aiServiceMutex sync.Mutex
	executionCount int64
	isExecuting    int32
	lastExecTime   time.Time
)

// acquireLock 是一个带超时控制的执行锁获取函数。
func (m *Manager) acquireLock(skipIntervalCheck bool) error {
	if atomic.LoadInt32(&isExecuting) == 1 {
		return ErrServiceBusy
	}

	if !skipIntervalCheck && time.Since(lastExecTime) < m.config.MinInterval {
		return ErrServiceBusy
	}

	lockAcquired := make(chan bool, 1)
	go func() {
		aiServiceMutex.Lock()
		lockAcquired <- true
	}()

	select {
	case <-lockAcquired:
		if atomic.LoadInt32(&isExecuting) == 1 {
			aiServiceMutex.Unlock()
			return ErrServiceBusy
		}

		atomic.StoreInt32(&isExecuting, 1)
		atomic.AddInt64(&executionCount, 1)
		lastExecTime = time.Now()
		return nil
	case <-time.After(m.config.MaxExecutionTime):
		return ErrTimeout
	}
}

// releaseLock 是对应的执行锁释放函数。
func (m *Manager) releaseLock() {
	atomic.StoreInt32(&isExecuting, 0)
	aiServiceMutex.Unlock()
}

// GetExecutionStatus 是一个状态查询函数，返回执行锁的状态数据。
func GetExecutionStatus() map[string]interface{} {
	return map[string]interface{}{
		"is_executing":    atomic.LoadInt32(&isExecuting) == 1,
		"execution_count": atomic.LoadInt64(&executionCount),
		"last_exec_time":  lastExecTime,
		"time_since_last": time.Since(lastExecTime),
	}
}

