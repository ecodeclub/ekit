package celery_pool

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type fastTask struct {
}

func (f *fastTask) Run(ctx context.Context) error {
	fmt.Println("This is fast Task")
	for i := 0; i < 3; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(1 * time.Second)

		}
	}
	return nil
}

type slowTask struct {
}

func (s *slowTask) Run(ctx context.Context) error {
	fmt.Println("This is slow Task")
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(1 * time.Second)

		}
	}
	return nil
}

func TestInMemoryTaskPool_Submit(t1 *testing.T) {

	t1.Run("short task", func(t *testing.T) {
		p := NewInMemoryTaskPool(WithConfig(poolConfig{
			checkDuration:        1 * time.Second,
			maxNormalWorker:      2,
			maxLongWorker:        1,
			normalExceedDuration: 5 * time.Second,
		}))
		task := fastTask{}
		cTask, err := p.Submit(context.Background(), &task)
		assert.NoError(t, err)
		startTime := time.Now()
		go p.Start()
		time.Sleep(4 * time.Second)
		p.Shutdown()
		assert.True(t, time.Since(cTask.GetStatus().pendingN) > time.Since(startTime))
		assert.False(t, time.Since(cTask.GetStatus().pendingL) > time.Since(startTime))
		assert.Equal(t, FINISH, cTask.GetStatus().state)
	})

	t1.Run("long task", func(t *testing.T) {
		p := NewInMemoryTaskPool(WithConfig(poolConfig{
			checkDuration:        1 * time.Second,
			maxNormalWorker:      2,
			maxLongWorker:        1,
			normalExceedDuration: 5 * time.Second,
		}))
		task := slowTask{}
		startTime := time.Now()

		cTask, err := p.Submit(context.Background(), &task)
		assert.NoError(t, err)
		go p.Start()
		ch, _ := p.Shutdown()
		for {
			select {
			case <-ch:
				assert.True(t, time.Since(cTask.GetStatus().pendingN) < time.Since(startTime))
				assert.True(t, time.Since(cTask.GetStatus().pendingL) < time.Since(startTime))
				assert.Equal(t, FINISH, cTask.GetStatus().state)
				return
			default:
				time.Sleep(1 * time.Second)
			}
		}

	})

	t1.Run("close now", func(t *testing.T) {
		p := NewInMemoryTaskPool(WithConfig(poolConfig{
			checkDuration:        1 * time.Second,
			maxNormalWorker:      2,
			maxLongWorker:        1,
			normalExceedDuration: 5 * time.Second,
		}))
		task := slowTask{}
		startTime := time.Now()

		cTask, err := p.Submit(context.Background(), &task)
		assert.NoError(t, err)
		go p.Start()
		time.Sleep(10 * time.Second)
		tasks, err := p.ShutdownNow()
		assert.NoError(t, err)
		assert.Contains(t, tasks, cTask)
		assert.True(t, time.Since(cTask.GetStatus().pendingN) < time.Since(startTime))
		assert.True(t, time.Since(cTask.GetStatus().pendingL) < time.Since(startTime))
		assert.Equal(t, ERROR, cTask.GetStatus().state)
		assert.Equal(t, context.Canceled, cTask.GetStatus().err)
		return

	})

	t1.Run("multiple tasks", func(t *testing.T) {
		p := NewInMemoryTaskPool(WithConfig(poolConfig{
			checkDuration:        1 * time.Second,
			maxNormalWorker:      5,
			maxLongWorker:        5,
			normalExceedDuration: 5 * time.Second,
		}))
		//startTime := time.Now()
		var cTasks []CeleryTask
		go p.Start()
		for i := 0; i < 10; i++ {
			s := slowTask{}
			cTask, err := p.Submit(context.Background(), &s)
			cTasks = append(cTasks, cTask)
			f := fastTask{}
			cTask, err = p.Submit(context.Background(), &f)
			assert.NoError(t, err)
			cTasks = append(cTasks, cTask)
		}

		ch, _ := p.Shutdown()
		for {
			select {
			case <-ch:
				for _, ct := range cTasks {
					assert.Equal(t, FINISH, ct.GetStatus().state)
				}
				return
			default:
				time.Sleep(1 * time.Second)
				for _, ct := range cTasks {
					t.Logf("%v", ct.GetStatus())
				}
			}
		}
		return

	})
}
