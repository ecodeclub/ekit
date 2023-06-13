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
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestCondSignal(t *testing.T) {
	var m sync.Mutex
	c := NewCond(&m)
	n := 2
	running := make(chan bool, n)
	awake := make(chan bool, n)
	for i := 0; i < n; i++ {
		go func() {
			m.Lock()
			running <- true
			c.Wait()
			awake <- true
			m.Unlock()
		}()
	}
	for i := 0; i < n; i++ {
		<-running // Wait for everyone to run.
	}
	for n > 0 {
		select {
		case <-awake:
			t.Fatal("goroutine not asleep")
		default:
		}
		m.Lock()
		c.Signal()
		m.Unlock()
		<-awake // Will deadlock if no goroutine wakes up
		select {
		case <-awake:
			t.Fatal("too many goroutines awake")
		default:
		}
		n--
	}
	c.Signal()
}

func TestCondSignalGenerations(t *testing.T) {
	var m sync.Mutex
	c := NewCond(&m)
	n := 100
	running := make(chan bool, n)
	awake := make(chan int, n)
	for i := 0; i < n; i++ {
		go func(i int) {
			m.Lock()
			running <- true
			c.Wait()
			awake <- i
			m.Unlock()
		}(i)
		if i > 0 {
			a := <-awake
			if a != i-1 {
				t.Fatalf("wrong goroutine woke up: want %d, got %d", i-1, a)
			}
		}
		<-running
		m.Lock()
		c.Signal()
		m.Unlock()
	}
}

func TestCondBroadcast(t *testing.T) {
	var m sync.Mutex
	c := NewCond(&m)
	n := 200
	running := make(chan int, n)
	awake := make(chan int, n)
	exit := false
	for i := 0; i < n; i++ {
		go func(g int) {
			m.Lock()
			for !exit {
				running <- g
				c.Wait()
				awake <- g
			}
			m.Unlock()
		}(i)
	}
	for i := 0; i < n; i++ {
		for i := 0; i < n; i++ {
			<-running // Will deadlock unless n are running.
		}
		if i == n-1 {
			m.Lock()
			exit = true
			m.Unlock()
		}
		select {
		case <-awake:
			t.Fatal("goroutine not asleep")
		default:
		}
		m.Lock()
		c.Broadcast()
		m.Unlock()
		seen := make([]bool, n)
		for i := 0; i < n; i++ {
			g := <-awake
			if seen[g] {
				t.Fatal("goroutine woke up twice")
			}
			seen[g] = true
		}
	}
	select {
	case <-running:
		t.Fatal("goroutine did not exit")
	default:
	}
	c.Broadcast()
}

func TestRace(t *testing.T) {
	x := 0
	c := NewCond(&sync.Mutex{})
	done := make(chan bool)
	go func() {
		c.L.Lock()
		x = 1
		c.Wait()
		if x != 2 {
			t.Error("want 2")
		}
		x = 3
		c.Signal()
		c.L.Unlock()
		done <- true
	}()
	go func() {
		c.L.Lock()
		for {
			if x == 1 {
				x = 2
				c.Signal()
				break
			}
			c.L.Unlock()
			runtime.Gosched()
			c.L.Lock()
		}
		c.L.Unlock()
		done <- true
	}()
	go func() {
		c.L.Lock()
		for {
			if x == 2 {
				c.Wait()
				if x != 3 {
					t.Error("want 3")
				}
				break
			}
			if x == 3 {
				break
			}
			c.L.Unlock()
			runtime.Gosched()
			c.L.Lock()
		}
		c.L.Unlock()
		done <- true
	}()
	<-done
	<-done
	<-done
}

func TestCondSignalStealing(t *testing.T) {
	for iters := 0; iters < 1000; iters++ {
		var m sync.Mutex
		cond := NewCond(&m)

		// Start a waiter.
		ch := make(chan struct{})
		go func() {
			m.Lock()
			ch <- struct{}{}
			cond.Wait()
			m.Unlock()

			ch <- struct{}{}
		}()

		<-ch
		m.Lock()
		done := false
		m.Unlock()

		// We know that the waiter is in the cond.Wait() call because we
		// synchronized with it, then acquired/released the mutex it was
		// holding when we synchronized.
		//
		// Start two goroutines that will race: one will broadcast on
		// the cond var, the other will wait on it.
		//
		// The new waiter may or may not get notified, but the first one
		// has to be notified.

		go func() {
			cond.Broadcast()
		}()

		go func() {
			m.Lock()
			for !done {
				cond.Wait()
			}
			m.Unlock()
		}()

		// Check that the first waiter does get signaled.
		select {
		case <-ch:
		case <-time.After(2 * time.Second):
			t.Fatalf("First waiter didn't get broadcast.")
		}

		// Release the second waiter in case it didn't get the
		// broadcast.
		m.Lock()
		done = true
		m.Unlock()
		cond.Broadcast()
	}
}

func BenchmarkCond1(b *testing.B) {
	benchmarkCond(b, 1)
}

func BenchmarkCond2(b *testing.B) {
	benchmarkCond(b, 2)
}

func BenchmarkCond4(b *testing.B) {
	benchmarkCond(b, 4)
}

func BenchmarkCond8(b *testing.B) {
	benchmarkCond(b, 8)
}

func BenchmarkCond16(b *testing.B) {
	benchmarkCond(b, 16)
}

func BenchmarkCond32(b *testing.B) {
	benchmarkCond(b, 32)
}

func benchmarkCond(b *testing.B, waiters int) {
	c := NewCond(&sync.Mutex{})
	done := make(chan bool)
	id := 0

	for routine := 0; routine < waiters+1; routine++ {
		go func() {
			for i := 0; i < b.N; i++ {
				c.L.Lock()
				if id == -1 {
					c.L.Unlock()
					break
				}
				id++
				if id == waiters+1 {
					id = 0
					c.Broadcast()
				} else {
					c.Wait()
				}
				c.L.Unlock()
			}
			c.L.Lock()
			id = -1
			c.Broadcast()
			c.L.Unlock()
			done <- true
		}()
	}
	for routine := 0; routine < waiters+1; routine++ {
		<-done
	}
}

func TestCond_WaitWithContext(t *testing.T) {

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
			err := cond.WaitWithContext(ctx)
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

func TestCond_WakeOrder(t *testing.T) {

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
			err := cond.WaitWithContext(ctx)
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
