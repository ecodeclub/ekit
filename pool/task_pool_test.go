package pool

import (
	"context"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// new an err pool
	_, e := NewBlockQueueTaskPool(1, -1)
	assert.Equal(t, errData, e)
	_, e = NewBlockQueueTaskPool(-3, 1)
	assert.Equal(t, errData, e)

	// new an opened pool
	_ = startTaskPool(1, 1, 1)
	r, e := NewBlockQueueTaskPool(1, 1)
	assert.Nil(t, e)
	assert.Equal(t, statusNew, r.showPoolStatus()) // check status

	// new a shutdown pool
	_ = startTaskPool(2, 2, 2)
	r, e = NewBlockQueueTaskPool(2, 2)
	assert.Nil(t, e)
	assert.Equal(t, statusNew, r.showPoolStatus())

	// new a shutdown_now pool
	_ = startTaskPool(3, 3, 3)
	r, e = NewBlockQueueTaskPool(3, 3)
	assert.Nil(t, e)
	assert.Equal(t, statusNew, r.showPoolStatus())

	// no-error start
	r, err := NewBlockQueueTaskPool(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, statusNew, r.showPoolStatus())
}

func TestBlockQueueTaskPool_Start(t *testing.T) {
	testCases := []struct {
		name       string
		list       *BlockQueueTaskPool
		wantStatus int
		wantErr    error
	}{
		{
			name:    "start",
			list:    startTaskPool(1, 1, 0),
			wantErr: nil,
		},
		{
			name:    "start opened",
			list:    startTaskPool(1, 1, 1),
			wantErr: errPoolOpened,
		},
		{
			name:    "start Shutdown",
			list:    startTaskPool(1, 1, 2),
			wantErr: errPoolStatus,
		},
		{
			name:    "start ShutdownNow",
			list:    startTaskPool(1, 1, 3),
			wantErr: errPoolStatus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Start()
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestBlockQueueTaskPool_Submit(t *testing.T) {
	// submit a nil-task
	p := startTaskPool(3, 3, 0)
	e := p.Submit(p.ctx, nil)
	assert.Equal(t, errTaskEmpty, e)

	// submit a task to new pool
	p = startTaskPool(3, 3, 0)
	e = p.Submit(p.ctx, TaskFunc(testSubmitFunc()))
	assert.Equal(t, errPoolStatus, e)

	// submit a task to shut down pool
	p = startTaskPool(3, 3, 2)
	e = p.Submit(p.ctx, TaskFunc(testSubmitFunc()))
	assert.Equal(t, errPoolStatus, e)

	// submit a task to shutdown-now pool
	p = startTaskPool(3, 3, 3)
	e = p.Submit(p.ctx, TaskFunc(testSubmitFunc()))
	assert.Equal(t, errPoolStatus, e)

	// submit a task to an opened pool
	p = startTaskPool(1, 1, 1)
	e = p.Submit(p.ctx, TaskFunc(testSubmitFunc()))
	assert.Nil(t, e)
	assert.Equal(t, statusOpen, p.showPoolStatus())

	// test cancel()
	p = startTaskPool(1, 1, 1)
	defer p.ctxCancel()
	for i := 0; i < 5; i++ {
		e = p.Submit(p.ctx, TaskFunc(testSubmitFunc()))
		assert.Nil(t, e)
		assert.Equal(t, statusOpen, p.showPoolStatus())
	}
}

func TestBlockQueueTaskPool_Shutdown(t *testing.T) {
	// shutdown a new pool
	p := startTaskPool(1, 1, 0)
	_, e := p.Shutdown()
	assert.Nil(t, e)
	assert.Equal(t, statusClose, p.showPoolStatus())

	// shutdown an opened pool
	p = startTaskPool(1, 1, 1)
	_, e = p.Shutdown()
	assert.Nil(t, e)
	assert.Equal(t, statusClose, p.showPoolStatus())

	// shutdown a shutdown-pool
	p = startTaskPool(2, 2, 2)
	_, e = p.Shutdown()
	assert.Equal(t, errPoolStatus, e)
	assert.Equal(t, statusClose, p.showPoolStatus())

	// shutdown a shutdown-now-pool
	p = startTaskPool(3, 3, 3)
	_, e = p.Shutdown()
	assert.Equal(t, errPoolStatus, e)
	assert.Equal(t, statusClose, p.showPoolStatus())
}

func TestBlockQueueTaskPool_ShutdownNow(t *testing.T) {
	// shutdown-now a new pool
	p := startTaskPool(1, 1, 0)
	_, e := p.ShutdownNow()
	assert.Nil(t, e)
	assert.Equal(t, statusClose, p.showPoolStatus())

	// shutdown-now an opened pool
	p = startTaskPool(1, 1, 1)
	_, e = p.ShutdownNow()
	assert.Nil(t, e)
	assert.Equal(t, statusClose, p.showPoolStatus())

	// shutdown-now a shutdown-pool
	p = startTaskPool(2, 2, 2)
	_, e = p.ShutdownNow()
	assert.Equal(t, errPoolStatus, e)
	assert.Equal(t, statusClose, p.showPoolStatus())

	// shutdown-now a shutdown-now-pool
	p = startTaskPool(3, 3, 3)
	_, e = p.ShutdownNow()
	assert.Equal(t, errPoolStatus, e)
	assert.Equal(t, statusClose, p.showPoolStatus())
}

func startTaskPool(concurrency, queueSize, flag int) *BlockQueueTaskPool {
	r, _ := NewBlockQueueTaskPool(concurrency, queueSize)

	switch flag {
	case 1:
		_ = r.Start()
	case 2:
		_ = r.Start()
		_, _ = r.Shutdown()
	case 3:
		_ = r.Start()
		_, _ = r.ShutdownNow()
	}
	return r
}

func testSubmitFunc() func(ctx context.Context) error {
	return func(ctx context.Context) error {
		time.Sleep(6 * time.Second)
		spew.Dump("a task finished")
		return nil
	}
}
