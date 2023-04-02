package celery_pool

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/pool"
	"sync"
	"time"
)

type TaskPool interface {
	Submit(ctx context.Context, task pool.Task) (CeleryTask, error)
	Start() error

	// Shutdown 关闭任务池。如果此时尚未调用 Start 方法，那么将会立刻返回。
	// 任务池将会停止接收新的任务，但是会继续执行剩下的任务，
	// 在所有任务执行完毕之后，用户可以从返回的 chan 中得到通知
	// 任务池在发出通知之后会关闭 chan struct{}
	Shutdown() (<-chan struct{}, error)

	// ShutdownNow 立刻关闭线程池
	// 任务池能否中断当前正在执行的任务，取决于 TaskPool 的具体实现，以及 Task 的具体实现
	// 该方法会返回所有剩下的任务，剩下的任务是否包含正在执行的任务，也取决于具体的实现
	ShutdownNow() ([]CeleryTask, error)
}

type InMemoryTaskPool struct {
	isClose    bool
	finishChan chan struct{}

	normalExceedDuration time.Duration
	longExceedDuration   time.Duration

	maxNormalWorker    int
	maxLongWorker      int
	checkDuration      time.Duration
	normalPendingQueue chan CeleryTask
	longPendingQueue   chan CeleryTask

	closeSignal chan struct{}
	taskMap     map[string]CeleryTask
}

type poolConfig struct {
	normalExceedDuration time.Duration
	longExceedDuration   time.Duration
	maxNormalWorker      int
	maxLongWorker        int
	checkDuration        time.Duration
}

type TaskPoolOption func(taskPool *InMemoryTaskPool)

func WithConfig(config poolConfig) TaskPoolOption {
	return func(taskPool *InMemoryTaskPool) {
		if config.normalExceedDuration > time.Second*0 {
			taskPool.normalExceedDuration = config.normalExceedDuration

		}
		if config.longExceedDuration > time.Second*0 {
			taskPool.longExceedDuration = config.longExceedDuration

		}
		if config.maxLongWorker > 0 {
			taskPool.maxLongWorker = config.maxLongWorker

		}
		if config.maxNormalWorker > 0 {
			taskPool.maxNormalWorker = config.maxNormalWorker

		}
		if config.checkDuration > time.Second*0 {

			taskPool.checkDuration = config.checkDuration
		}
	}
}

func NewInMemoryTaskPool(opts ...TaskPoolOption) *InMemoryTaskPool {
	res := &InMemoryTaskPool{
		isClose:              false,
		finishChan:           make(chan struct{}, 1),
		normalExceedDuration: time.Minute,
		longExceedDuration:   time.Hour,

		maxNormalWorker: 100,
		maxLongWorker:   10,
		checkDuration:   time.Second,
		taskMap:         make(map[string]CeleryTask),
		closeSignal:     make(chan struct{}),
	}
	for _, opt := range opts {
		opt(res)
	}
	res.normalPendingQueue = make(chan CeleryTask, res.maxNormalWorker)
	res.longPendingQueue = make(chan CeleryTask, res.maxLongWorker)
	return res
}

func (t *InMemoryTaskPool) Submit(ctx context.Context, task pool.Task) (CeleryTask, error) {
	if t.isClose {
		return nil, errors.New("pool already close")
	}
	ct := &InMemoryTask{
		task: task.Run,
	}
	t.taskMap[ct.GetId()] = ct
	ct.setState(PENDING_N, nil)
	t.normalPendingQueue <- ct
	return ct, nil
}

func (t *InMemoryTaskPool) Start() error {
	wg := sync.WaitGroup{}
	for i := 0; i < t.maxNormalWorker; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case task := <-t.normalPendingQueue:
					ctx, cancel := context.WithTimeout(context.Background(), t.normalExceedDuration)
					task.setState(START_N, nil)
					task.setCancel(cancel)
					go func() {
						err := task.Run(ctx)
						defer cancel()
						if err != nil {
							//delete(t.taskMap, task.GetId())
							task.setState(ERROR, err)
							return
						}
						task.setState(FINISH, nil)
					}()

					select {
					case <-ctx.Done():
						err := ctx.Err()
						if err == context.DeadlineExceeded {
							task.setState(PENDING_L, nil)
							task.setCancel(nil)
							t.longPendingQueue <- task
						} else {
							delete(t.taskMap, task.GetId())
						}
					}
				case <-t.closeSignal:
					break
				}
			}
			wg.Done()
		}()
	}
	for i := 0; i < t.maxLongWorker; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case task := <-t.longPendingQueue:
					ctx, cancel := context.WithTimeout(context.Background(), t.longExceedDuration)
					task.setState(START_N, nil)
					task.setCancel(cancel)
					go func() {
						err := task.Run(ctx)
						defer cancel()
						if err != nil {
							//delete(t.taskMap, task.GetId())
							task.setState(ERROR, err)
							return
						}
						task.setState(FINISH, nil)
					}()

					select {
					case <-ctx.Done():
						err := ctx.Err()
						if err == context.DeadlineExceeded {
							task.setState(ERROR, err)
						}
						delete(t.taskMap, task.GetId())

					}
				case <-t.closeSignal:
					break
				}
			}
			wg.Done()
		}()

	}

	ticker := time.NewTicker(t.checkDuration)
	for {
		select {
		case <-ticker.C:
			if len(t.taskMap) == 0 && t.isClose {
				t.finishChan <- struct{}{}
				t.closeSignal <- struct{}{}
			}
		case <-t.closeSignal:
			wg.Wait()
			t.closeChan()
			return nil
		}
	}
}

func (t *InMemoryTaskPool) Shutdown() (<-chan struct{}, error) {
	t.isClose = true
	return t.finishChan, nil

}

func (t *InMemoryTaskPool) ShutdownNow() ([]CeleryTask, error) {
	t.isClose = true
	pendingTasks := make([]CeleryTask, len(t.taskMap))
	for _, task := range t.taskMap {
		pendingTasks = append(pendingTasks, task)
		task.Cancel()
	}
	t.closeSignal <- struct{}{}
	t.closeSignal <- struct{}{}
	return pendingTasks, nil
}

func (t *InMemoryTaskPool) closeChan() {
	close(t.closeSignal)
	close(t.longPendingQueue)
	close(t.normalPendingQueue)
	close(t.finishChan)
}
