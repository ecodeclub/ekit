package pool

import (
	"context"
	"time"
)

type MetricTask struct {
	Task
	observer TaskObserver
	state    int32
	logFunc  func(args ...any)
}

// InternalState 用于查看Task状态
func (t *MetricTask) InternalState() int32 {
	panic("Task State")
}

func (t *MetricTask) Run(ctx context.Context) error {
	startTime := time.Now().UnixNano()
	err := t.Task.Run(ctx)
	t.logFunc(err)
	endTime := time.Now().UnixNano()
	t.observer.Observe(ctx, t, endTime-startTime)
	return err
}

type MetricTaskPool struct {
	TaskPool
	observer TaskPoolObserver
	logFunc  func(args ...any)
}

func (p *MetricTaskPool) Start() error {
	startTime := time.Now().UnixNano()
	err := p.TaskPool.Start()
	p.logFunc(err)
	endTime := time.Now().UnixNano()
	p.observer.Observe(context.Background(), p, endTime-startTime)
	return err
}
