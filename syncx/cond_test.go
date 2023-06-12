package syncx

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
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
		done := false
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

func TestNotifyListSignal(t *testing.T) {
	nl := newNotifyList()

	wait1 := nl.add()
	wait2 := nl.add()
	wait3 := nl.add()
	wait4 := nl.add()
	wait5 := nl.add()
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		wg.Done()
		wait1.wait()
	}()
	go func() {
		wg.Done()
		wait2.wait()
	}()
	go func() {
		wg.Done()
		wait3.wait()
	}()
	go func() {
		wg.Done()
		wait4.wait()
	}()
	go func() {
		wg.Done()
		wait5.wait()
	}()
	wg.Wait()
	nl.notifyOne()
	wait1.wait()
	nl.notifyOne()
	wait2.wait()
	nl.notifyOne()
	wait3.wait()
	nl.notifyOne()
	wait4.wait()
	nl.notifyOne()
	wait5.wait()
}

func TestNotifyListBroadcast(t *testing.T) {
	nl := newNotifyList()

	wait1 := nl.add()
	wait2 := nl.add()
	wait3 := nl.add()
	wait4 := nl.add()
	wait5 := nl.add()
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		wg.Done()
		wait1.wait()
	}()
	go func() {
		wg.Done()
		wait2.wait()
	}()
	go func() {
		wg.Done()
		wait3.wait()
	}()
	go func() {
		wg.Done()
		wait4.wait()
	}()
	go func() {
		wg.Done()
		wait5.wait()
	}()
	wg.Wait()
	nl.notifyAll()
	wait1.wait()
	wait2.wait()
	wait3.wait()
	wait4.wait()
	wait5.wait()
}

func TestNotifyListSignalTimeout(t *testing.T) {
	nl := newNotifyList()

	wait1 := nl.add()
	wait2 := nl.add()
	wait3 := nl.add()
	wait4 := nl.add()
	wait5 := nl.add()
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		wg.Done()
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancelFunc()
		err := wait1.waitWithContext(ctx)
		if err != nil {

		}
	}()
	go func() {
		wg.Done()
		wait2.wait()
	}()
	go func() {
		wg.Done()
		wait3.wait()
	}()
	go func() {
		wg.Done()
		wait4.wait()
	}()
	go func() {
		wg.Done()
		wait5.wait()
	}()
	wg.Wait()
	time.Sleep(time.Millisecond * 150)
	nl.notifyOne()
	wait2.wait()
	nl.notifyOne()
	wait3.wait()
	nl.notifyOne()
	wait4.wait()
	nl.notifyOne()
	wait5.wait()
}

func TestCond(t *testing.T) {
	x := 0
	c := NewCond(&sync.Mutex{})
	done := make(chan bool)
	go func() {
		c.L.Lock()
		x = 1
		c.Wait()
		if x != 2 {
			log.Fatal("want 2")
		}
		x = 3
		c.Broadcast()
		c.L.Unlock()
		done <- true
	}()
	go func() {
		c.L.Lock()
		for {
			if x == 1 {
				x = 2
				c.Broadcast()
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
					log.Fatal("want 3")
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

// 使用并发阻塞队列进行测试生产和消费的情况
type ConcurrentBlockingQueue[T any] struct {
	mutex *sync.Mutex
	data  []T
	// notFull chan struct{}
	// notEmpty chan struct{}
	maxSize int

	notEmptyCond *Cond
	notFullCond  *Cond
}

func NewConcurrentBlockingQueue[T any](maxSize int) *ConcurrentBlockingQueue[T] {
	m := &sync.Mutex{}
	return &ConcurrentBlockingQueue[T]{
		data:  make([]T, 0, maxSize),
		mutex: m,
		// notFull: make(chan struct{}, 1),
		// notEmpty: make(chan struct{}, 1),
		maxSize:      maxSize,
		notFullCond:  NewCond(m),
		notEmptyCond: NewCond(m),
	}
}

func (c *ConcurrentBlockingQueue[T]) EnQueue(ctx context.Context, data T) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for c.isFull() {
		err := c.notFullCond.WaitWithContext(ctx)
		if err != nil {
			return err
		}
	}
	c.data = append(c.data, data)
	c.notEmptyCond.Signal()
	// 没有人等 notEmpty 的信号，这一句就会阻塞住
	return nil
}

func (c *ConcurrentBlockingQueue[T]) DeQueue(ctx context.Context) (T, error) {
	if ctx.Err() != nil {
		var t T
		return t, ctx.Err()
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for c.isEmpty() {
		err := c.notEmptyCond.WaitWithContext(ctx)
		if err != nil {
			var t T
			return t, err
		}
	}
	t := c.data[0]
	c.data = c.data[1:]
	c.notFullCond.Signal()
	// 没有人等 notFull 的信号，这一句就会阻塞住
	return t, nil
}

func (c *ConcurrentBlockingQueue[T]) IsFull() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.isFull()
}

func (c *ConcurrentBlockingQueue[T]) isFull() bool {
	return len(c.data) == c.maxSize
}

func (c *ConcurrentBlockingQueue[T]) IsEmpty() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.isEmpty()
}

func (c *ConcurrentBlockingQueue[T]) isEmpty() bool {
	return len(c.data) == 0
}

func (c *ConcurrentBlockingQueue[T]) Len() uint64 {
	return uint64(len(c.data))
}

func TestConcurrentBlockingQueue_EnQueue(t *testing.T) {
	testCases := []struct {
		name string

		q *ConcurrentBlockingQueue[int]

		timeout time.Duration
		value   int

		data []int

		wantErr error
	}{
		{
			name:    "enqueue",
			q:       NewConcurrentBlockingQueue[int](10),
			value:   1,
			timeout: time.Minute,
			data:    []int{1},
		},
		{
			name: "blocking and timeout",
			q: func() *ConcurrentBlockingQueue[int] {
				res := NewConcurrentBlockingQueue[int](2)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := res.EnQueue(ctx, 1)
				require.NoError(t, err)
				err = res.EnQueue(ctx, 2)
				require.NoError(t, err)
				return res
			}(),
			value:   3,
			timeout: time.Second,
			data:    []int{1, 2},
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()
			err := tc.q.EnQueue(ctx, tc.value)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.data, tc.q.data)
		})
	}
}

// 测试可能比较花时间，不知道如何优化这个模拟延迟的缓慢问题
func TestConcurrentBlockingQueue(t *testing.T) {
	q := NewConcurrentBlockingQueue[int](20)
	var wg sync.WaitGroup
	var producedCnt, consumedCnt int64

	producers := [20]int{}
	consumers := [10]int{}
	// 使用2w个数据进行测试，使用随机睡眠模拟生产者和消费者
	go func() {
		for range time.Tick(time.Second) {
			// 每秒打印生产和消费的线程状态和成功生产和消费总数
			log.Println(producers, consumers, atomic.LoadInt64(&producedCnt), atomic.LoadInt64(&consumedCnt))
		}
	}()
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(idx int) {
			defer func() {
				producers[idx] = 1
				wg.Done()
			}()
			n := 1000
			for n > 0 {
				n--
				// 随机睡眠模拟生产者 -- 模拟出现消费者将队列消费空之后出现error的情况
				time.Sleep(time.Duration(rand.Int()%20) * time.Millisecond)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
				err := q.EnQueue(ctx, rand.Int())
				if err == nil {
					// 成功生产的数据量
					atomic.AddInt64(&producedCnt, 1)
				}
				//fmt.Println(err)
				// 怎么断言 error
				cancel()
			}
		}(i)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer func() {
				consumers[idx] = 1
				wg.Done()
			}()
			n := 2000
			for n > 0 {
				n--
				// 随机睡眠模拟消费者 -- 模拟出现生产者将队列放满之后出现error的情况
				time.Sleep(time.Duration(rand.Int()%20) * time.Millisecond)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
				_, err := q.DeQueue(ctx)
				if err == nil {
					// 成功消费的数据量
					atomic.AddInt64(&consumedCnt, 1)
				}
				//fmt.Println(val, err)
				// 又怎么断言 val, 和 err
				cancel()
			}
		}(i)
	}
	wg.Wait()
	fmt.Println(producedCnt, consumedCnt)
	// 怎么校验 q 对还是不对
}
