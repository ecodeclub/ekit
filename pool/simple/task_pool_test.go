// Copyright 2021 gotomicro
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pool

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBlockQueueTaskPool_Start(t *testing.T) {
	testCases := []struct {
		name        string
		instance    *BlockQueueTaskPool
		expectedErr error
	}{
		{
			name:        "正常启动",
			instance:    newTaskPool(2, 4, false),
			expectedErr: nil,
		},
		{
			name:        "重复启动",
			instance:    newTaskPool(2, 4, true),
			expectedErr: errTaskPoolAlreadyStarted,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.instance.Start()
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestBlockQueueTaskPool_Submit(t *testing.T) {
	testCases := []struct {
		name        string
		instance    *BlockQueueTaskPool
		task        TaskFunc
		expectedErr error
	}{
		{
			name:        "向未启动的任务池提交任务",
			instance:    newTaskPool(3, 3, false),
			task:        createF(),
			expectedErr: nil,
		},
		{
			name:        "向启动启动的任务池提交任务",
			instance:    newTaskPool(3, 3, true),
			task:        createF(),
			expectedErr: nil,
		},
		{
			name:        "向满队的任务池提交任务",
			instance:    newTaskPool(0, 0, true),
			task:        createF(),
			expectedErr: context.DeadlineExceeded,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.instance.Submit(ctx, tc.task)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestBlockQueueTaskPool_Shutdown(t *testing.T) {
	testCases := []struct {
		name        string
		instance    *BlockQueueTaskPool
		closeSignal struct{}
		expectedErr error
	}{
		{
			name:        "关闭已启动的空任务池",
			instance:    newTaskPoolWithStatus(2, 2, true, 0),
			closeSignal: struct{}{},
			expectedErr: nil,
		},
		{
			name:        "关闭已启动的非空任务池",
			instance:    newTaskPoolWithStatus(1, 1, true, 3),
			closeSignal: struct{}{},
			expectedErr: nil,
		},
		{
			name:        "关闭未启动的空任务池",
			instance:    newTaskPoolWithStatus(2, 2, false, 0),
			closeSignal: struct{}{},
			expectedErr: nil,
		},
		{
			name:        "关闭未启动的非空任务池",
			instance:    newTaskPoolWithStatus(2, 2, false, 3),
			closeSignal: struct{}{},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			signalChan, err := tc.instance.Shutdown()
			assert.Equal(t, tc.expectedErr, err)
			if err != nil {
				return
			}
			select {
			case signal := <-signalChan:
				assert.Equal(t, tc.closeSignal, signal)
			case <-time.After(time.Minute * 2):
				t.Fatal("未正常关闭")
			}
		})
	}
}
func TestBlockQueueTaskPool_ShutdownMulti(t *testing.T) {
	testCases := []struct {
		name        string
		instance    *BlockQueueTaskPool
		closeSignal struct{}
		expectedErr error
	}{
		{
			name:        "重复关闭已启动的空任务池",
			instance:    newTaskPoolWithStatus(2, 2, true, 0),
			closeSignal: struct{}{},
			expectedErr: errTaskPoolClosed,
		},
		{
			name:        "重复关闭已启动的非空任务池",
			instance:    newTaskPoolWithStatus(1, 1, true, 3),
			closeSignal: struct{}{},
			expectedErr: errTaskPoolClosed,
		},
		{
			name:        "重复关闭未启动的空任务池",
			instance:    newTaskPoolWithStatus(2, 2, false, 0),
			closeSignal: struct{}{},
			expectedErr: errTaskPoolClosed,
		},
		{
			name:        "重复关闭未启动的非空任务池",
			instance:    newTaskPoolWithStatus(2, 2, false, 3),
			closeSignal: struct{}{},
			expectedErr: errTaskPoolClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.instance.Shutdown()
			if err != nil {
				return
			}
			_, err = tc.instance.Shutdown()
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestBlockQueueTaskPool_ShutdownNow(t *testing.T) {
	testCases := []struct {
		name          string
		instance      *BlockQueueTaskPool
		expectedTasks int
		expectedErr   error
	}{
		{
			name:          "立刻关闭已启动的空任务池",
			instance:      newTaskPoolWithStatus(2, 2, true, 0),
			expectedTasks: 0,
			expectedErr:   nil,
		},
		{
			name:          "立刻关闭已启动的非空任务池",
			instance:      newTaskPoolWithStatus(2, 10, true, 5),
			expectedTasks: 3,
			expectedErr:   nil,
		},
		{
			name:          "立刻关闭未启动的空任务池",
			instance:      newTaskPoolWithStatus(2, 2, false, 0),
			expectedTasks: 0,
			expectedErr:   nil,
		},
		{
			name:          "立刻关闭未启动的非空任务池",
			instance:      newTaskPoolWithStatus(2, 3, false, 3),
			expectedTasks: 3,
			expectedErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tasks, err := tc.instance.ShutdownNow()
			assert.Equal(t, tc.expectedErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.expectedTasks, len(tasks))
		})
	}
}

func TestBlockQueueTaskPool_ShutdownNowMulti(t *testing.T) {
	testCases := []struct {
		name          string
		instance      *BlockQueueTaskPool
		expectedTasks int
		expectedErr   error
	}{
		{
			name:        "立刻重复关闭已启动的空任务池",
			instance:    newTaskPoolWithStatus(2, 2, true, 0),
			expectedErr: errTaskPoolClosed,
		},
		{
			name:        "立刻重复关闭已启动的非空任务池",
			instance:    newTaskPoolWithStatus(2, 10, true, 5),
			expectedErr: errTaskPoolClosed,
		},
		{
			name:        "立刻重复关闭未启动的空任务池",
			instance:    newTaskPoolWithStatus(2, 2, false, 0),
			expectedErr: errTaskPoolClosed,
		},
		{
			name:        "立刻重复关闭未启动的非空任务池",
			instance:    newTaskPoolWithStatus(2, 3, false, 3),
			expectedErr: errTaskPoolClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.instance.ShutdownNow()
			if err != nil {
				return
			}
			_, err = tc.instance.ShutdownNow()
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestBlockQueueTaskPool_ShutdownAndShutdownNow(t *testing.T) {
	testCases := []struct {
		name        string
		instance    *BlockQueueTaskPool
		expectedErr error
	}{
		{
			name:        "关闭已启动的空任务池",
			instance:    newTaskPoolWithStatus(2, 2, true, 0),
			expectedErr: errTaskPoolClosed,
		},
		{
			name:        "关闭已启动的非空任务池",
			instance:    newTaskPoolWithStatus(2, 10, true, 5),
			expectedErr: errTaskPoolClosed,
		},
		{
			name:        "关闭未启动的空任务池",
			instance:    newTaskPoolWithStatus(2, 2, false, 0),
			expectedErr: errTaskPoolClosed,
		},
		{
			name:        "关闭未启动的非空任务池",
			instance:    newTaskPoolWithStatus(2, 3, false, 3),
			expectedErr: errTaskPoolClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.instance.Shutdown()
			if err != nil {
				return
			}
			_, err = tc.instance.ShutdownNow()
			assert.Equal(t, tc.expectedErr, err)
		})
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.instance.ShutdownNow()
			if err != nil {
				return
			}
			_, err = tc.instance.Shutdown()
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestBlockQueueTaskPool_Submit_ShutdownPool(t *testing.T) {
	testCases := []struct {
		name        string
		instance    *BlockQueueTaskPool
		task        TaskFunc
		expectedErr error
	}{
		{
			name:        "shutdown任务池提交任务",
			instance:    newShutdownTaskPool(3, 3, false),
			task:        createF(),
			expectedErr: errTaskPoolClosed,
		},
		{
			name:        "shutdownNow的任务池提交任务",
			instance:    newShutdownTaskPool(3, 3, true),
			task:        createF(),
			expectedErr: errTaskPoolClosed,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.instance.Submit(ctx, tc.task)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func newTaskPool(concurrency, size int, started bool) *BlockQueueTaskPool {
	pool, _ := NewBlockQueueTaskPool(concurrency, size)
	if started {
		_ = pool.Start()
		time.Sleep(time.Second)
	}
	return pool
}

func newTaskPoolWithStatus(concurrency, size int, started bool, funcSize int) *BlockQueueTaskPool {
	pool := newTaskPool(concurrency, size, started)
	f := createF()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for i := 0; i < funcSize; i++ {
		_ = pool.Submit(ctx, f)
	}

	return pool
}

func newShutdownTaskPool(concurrency, size int, now bool) *BlockQueueTaskPool {
	pool, _ := NewBlockQueueTaskPool(concurrency, size)
	if now {
		_, _ = pool.ShutdownNow()
	} else {
		_, _ = pool.Shutdown()
	}
	time.Sleep(time.Second)
	return pool
}

func createF() TaskFunc {
	return func(ctx context.Context) error {
		time.Sleep(time.Millisecond * 10)
		fmt.Println("work is done")
		return nil
	}
}
