package pool

import (
	"context"
	"errors"
	"sync"
)

var (
	// 检测是否实现所有方法
	_ TaskPool = &BlockQueueTaskPool{}

	// 任务池中的状态码 1准备 2开启 3关闭
	statusNew   = 1
	statusOpen  = 2
	statusClose = 3

	// poolClosed 任务池已经关闭Error
	errData = errors.New("task pool:数据出错")
	// taskEmpty task为nil
	errTaskEmpty = errors.New("task pool:一个无效的task")
	// poolClosed 任务池已经关闭
	//errPoolClosed = errors.New("task pool:连接池已经关闭")
	// errPoolStatus 任务池状态错误
	errPoolStatus = errors.New("task pool:任务池状态错误")
	// errPoolOpened 任务池已经启动
	errPoolOpened = errors.New("task pool:任务池已经启动")
)

// TaskPool 任务池
type TaskPool interface {
	// Submit 执行一个任务
	// 如果任务池提供了阻塞的功能，那么如果在 ctx 过期都没有提交成功，那么应该返回错误
	// 调用 Start 之后能否继续提交任务，则取决于具体的实现
	// 调用 Shutdown 或者 ShutdownNow 之后提交任务都会返回错误
	Submit(ctx context.Context, task Task) error

	// Start 开始调度任务执行。在调用 Start 之前，所有的任务都不会被调度执行。
	// Start 之后，能否继续调用 Submit 提交任务，取决于具体的实现
	Start() error

	// Shutdown 关闭任务池。如果此时尚未调用 Start 方法，那么将会立刻返回。
	// 任务池将会停止接收新的任务，但是会继续执行剩下的任务，
	// 在所有任务执行完毕之后，用户可以从返回的 chan 中得到通知
	// 任务池在发出通知之后会关闭 chan struct{}
	Shutdown() (<-chan struct{}, error)

	// ShutdownNow 立刻关闭线程池
	// 任务池能否中断当前正在执行的任务，取决于 TaskPool 的具体实现，以及 Task 的具体实现
	// 该方法会返回所有剩下的任务，剩下的任务是否包含正在执行的任务，也取决于具体的实现
	ShutdownNow() ([]Task, error)
}

// Task 代表一个任务
type Task interface {
	// Run 执行任务
	// 如果 ctx 设置了超时时间，那么实现者需要自己决定是否进行超时控制
	Run(ctx context.Context) error
}

// BlockQueueTaskPool 并发阻塞的任务池
type BlockQueueTaskPool struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	waitTask  chan Task     // 等待的任务list
	doingTask chan struct{} // 正在执行的任务
	status    int           // 任务池的状态码
	mu        sync.Mutex
}
