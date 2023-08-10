// Copyright 2021 ecodeclub
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

package syncx

import (
	"context"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestCond_Broadcast(t *testing.T) {

	cond := NewCond(&sync.Mutex{})

	type status struct {
		i   int
		err error
	}

	sleepDuration := time.Millisecond * 100

	var n = 100
	running := make(chan int, n)
	awake := make(chan status, n)
	waitSeqs := make([]int, n)
	normalAwakeSeqs := make([]int, 0, n)
	timeoutAwakeSeqs := make([]int, 0, n)
	minTimeoutCnt := 0
	minNormalCnt := 0
	seen := make(map[int]bool, n)
	for i := 0; i < n; i++ {
		duration := time.Millisecond * 50 * time.Duration(rand.Int()%4+1)
		if duration < sleepDuration*9/10 {
			minTimeoutCnt++
		} else if duration > sleepDuration*11/10 {
			minNormalCnt++
		}
		go func(i int) {
			cond.L.Lock()

			ctx, cancelFunc := context.WithTimeout(context.Background(), duration)
			defer cancelFunc()
			running <- i
			err := cond.Wait(ctx)
			awake <- status{
				i:   i,
				err: err,
			}
			cond.L.Unlock()
		}(i)
	}
	for i := 0; i < n; i++ {
		waitSeqs[i] = <-running
	}

	time.Sleep(100 * time.Millisecond)

	cond.L.Lock()
	cond.Broadcast()
	cond.L.Unlock()

	for i := 0; i < n; i++ {
		stat := <-awake
		if seen[stat.i] {
			t.Fatal("goroutine woke up twice")
		} else {
			seen[stat.i] = true
		}
		if stat.err != nil {
			timeoutAwakeSeqs = append(timeoutAwakeSeqs, stat.i)
		} else {
			normalAwakeSeqs = append(normalAwakeSeqs, stat.i)
		}
	}

	if len(normalAwakeSeqs) < minNormalCnt {
		t.Fatal("goroutine woke up with timeout")
	}

	if len(timeoutAwakeSeqs) < minTimeoutCnt {
		t.Fatal("goroutine woke up with normally")
	}
}

func TestCond_Signal(t *testing.T) {

	cond := NewCond(&sync.Mutex{})

	type status struct {
		i   int
		err error
	}

	sleepDuration := time.Millisecond * 100

	var n = 100
	running := make(chan int, n)
	awake := make(chan status, n)
	waitSeqs := make([]int, n)
	normalAwakeSeqs := make([]int, 0, n)
	timeoutAwakeSeqs := make([]int, 0, n)
	minTimeoutCnt := 0
	minNormalCnt := 0
	seen := make(map[int]bool, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		duration := time.Millisecond * 50 * time.Duration(rand.Int()%4+1)
		if duration < sleepDuration*9/10 {
			minTimeoutCnt++
		} else if duration > sleepDuration*11/10 {
			minNormalCnt++
		}
		go func(i int) {
			cond.L.Lock()

			ctx, cancelFunc := context.WithTimeout(context.Background(), duration)
			defer cancelFunc()
			running <- i
			err := cond.Wait(ctx)
			awake <- status{
				i:   i,
				err: err,
			}
			cond.L.Unlock()
			wg.Done()
		}(i)
	}
	for i := 0; i < n; i++ {
		waitSeqs[i] = <-running
	}

	go func() {
		wg.Wait()
		close(awake)
	}()

	time.Sleep(100 * time.Millisecond)

	for i := 0; i < n; i++ {
		cond.L.Lock()
		cond.Signal()
		cond.L.Unlock()
		for {
			stat, ok := <-awake
			if !ok {
				break
			}
			if seen[stat.i] {
				t.Fatal("goroutine woke up twice")
			} else {
				seen[stat.i] = true
			}
			if stat.err != nil {
				timeoutAwakeSeqs = append(timeoutAwakeSeqs, stat.i)
			} else {
				normalAwakeSeqs = append(normalAwakeSeqs, stat.i)
				break
			}
		}

	}

	if len(normalAwakeSeqs) < minNormalCnt {
		t.Fatal("goroutine woke up with timeout")
	}

	if len(timeoutAwakeSeqs) < minTimeoutCnt {
		t.Fatal("goroutine woke up with normally")
	}
	// 测试singnal唤醒的顺序问题
	if !isInOrder(normalAwakeSeqs, waitSeqs) {
		t.Fatal("goroutine woke up not in order")
	}
	// 超时唤醒的肯定是乱序的，没有好办法测试顺序
	//if !isInOrder(timeoutAwakeSeqs, waitSeqs) {
	//	t.Fatal("goroutine woke up not in order")
	//}
}

func isInOrder(partial []int, source []int) bool {

	j := 0

	for i := 0; i < len(partial); i++ {
		matched := false
		for j < len(source) {
			if partial[i] == source[j] {
				j++
				matched = true
				break
			}
			j++
			continue
		}
		if !matched {
			return false
		}
	}

	return true
}

func Test_InOrder(t *testing.T) {
	testcases := []struct {
		name    string
		partial []int
		source  []int
		want    bool
	}{
		{"", []int{1}, []int{1}, true},
		{"", []int{1, 3, 4}, []int{1, 2, 3, 4}, true},
		{"", []int{1, 3}, []int{1, 2, 3, 4}, true},
		{"", []int{1, 3, 2}, []int{1, 2, 3, 4}, false},
		{"", []int{1, 2, 2}, []int{1, 2, 3, 4}, false},
		{"", []int{1, 2, 3}, []int{1, 3, 2, 4}, false},
		{"", []int{1, 2, 4}, []int{1, 3, 2, 4}, true},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			if target := isInOrder(tt.partial, tt.source); target != tt.want {
				t.Errorf("get %v, want %v", target, tt.want)
			}
		})
	}
}

// TestChanList 测试有序，和清空后重复使用是否有问题
func TestChanList(t *testing.T) {

	l := newChanList()

	testcases := []struct {
		name string
		num  int
	}{
		{"", 5},
		{"", 3},
		{"", 10},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(tt *testing.T) {
			inputNodes := make([]*node, 0, testcase.num)
			inputChans := make([]chan struct{}, 0, testcase.num)
			for i := 0; i < testcase.num; i++ {
				ele := l.alloc()
				inputNodes = append(inputNodes, ele)
				inputChans = append(inputChans, ele.Value)
				l.pushBack(ele)
			}
			if length := l.len(); length != testcase.num {
				t.Errorf("list.len() = %v, want %v", length, testcase.num)
			}
			outNodes := make([]*node, 0, testcase.num)
			outChans := make([]chan struct{}, 0, testcase.num)
			for l.len() != 0 {
				front := l.front()
				outNodes = append(outNodes, front)
				outChans = append(outChans, front.Value)
				l.remove(front)
			}
			if !reflect.DeepEqual(outChans, inputChans) {
				t.Errorf("chan list is %v, but got %v", inputChans, outChans)
			}
			if !reflect.DeepEqual(outNodes, inputNodes) {
				t.Errorf("element list is %v, but got %v", inputNodes, outNodes)
			}
		})
	}
}

// BenchmarkChanList 测试有无内存分配增加的情况
func BenchmarkChanList(b *testing.B) {
	l := newChanList()
	for i := 0; i < b.N; i++ {
		elem := l.alloc()
		l.pushBack(elem)
		l.remove(elem)
	}
}
