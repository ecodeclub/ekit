package celery_pool

import (
	"context"
	"github.com/ecodeclub/ekit/pool"
	"github.com/google/uuid"
	"time"
)

type TaskState string

const (
	PENDING_N TaskState = "pending normal"
	START_N   TaskState = "start normal"
	FINISH    TaskState = "finish"
	PENDING_L TaskState = "pending long"
	START_L   TaskState = "start long"
	RETRY     TaskState = "retry"
	ERROR     TaskState = "error"
)

type RunTimeStatus struct {
	pendingN   time.Time
	pendingL   time.Time
	startN     time.Time
	startL     time.Time
	finishTime time.Time
	state      TaskState
	err        error
}

func (rts *RunTimeStatus) Expired(expireDuration time.Duration) bool {
	return time.Since(rts.startN) > expireDuration
}

type CeleryTask interface {
	pool.Task
	setState(s TaskState, err error) error
	GetStatus() RunTimeStatus
	setCancel(c func())
	GetId() string
	Cancel() error
}

type InMemoryTask struct {
	uuid   string
	task   func(ctx context.Context) error
	status RunTimeStatus
	cancel func()
}

func (i *InMemoryTask) setCancel(c func()) {
	i.cancel = c
}

func (i *InMemoryTask) GetStatus() RunTimeStatus {
	return i.status
}

func (i *InMemoryTask) GetId() string {
	if i.uuid == "" {
		i.uuid = uuid.New().String()
	}
	return i.uuid
}

func (i *InMemoryTask) setState(s TaskState, err error) error {
	switch s {
	case START_N:
		i.status.state = START_N
		i.status.startN = time.Now()
	case START_L:
		i.status.state = START_L
		i.status.startL = time.Now()
	case PENDING_N:
		i.status.state = PENDING_N
		i.status.pendingN = time.Now()
	case PENDING_L:
		i.status.state = PENDING_L
		i.status.pendingL = time.Now()
	case FINISH:
		i.status.state = FINISH
		i.status.finishTime = time.Now()
	case ERROR:
		i.status.state = ERROR
		i.status.err = err
		i.status.finishTime = time.Now()
	}
	return nil
}

func (i *InMemoryTask) Run(ctx context.Context) error {
	return i.task(ctx)
}

func (i *InMemoryTask) Cancel() error {
	if i.cancel != nil {
		i.cancel()
		i.status.state = ERROR
		i.status.err = context.Canceled
	}

	return nil
}
