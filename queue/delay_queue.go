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

package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gotomicro/ekit"
	"github.com/gotomicro/ekit/internal/queue"
)

var (
	errInvalidArgument    = errors.New("ekit: 参数非法")
	errQueueHasBeenClosed = errors.New("ekit: 队列已关闭")
)

// Delayable 入队元素需要实现接口中的 Deadline 方法
type Delayable[T any] interface {
	Deadline() time.Time
}

// DelayQueue 并发延迟队列
type DelayQueue[T Delayable[T]] struct {
	mutex                *sync.RWMutex
	q                    *queue.PriorityQueue[T]
	capacity             int
	compareFuncOfElement ekit.Comparator[T]

	// Proxy 协程生命周期管理相关
	startedProxiesWaitGroup *sync.WaitGroup
	stoppedProxiesWaitGroup *sync.WaitGroup
	quitSignalChan          chan struct{}
	// enqueueProxy 协程相关
	numOfEnqueueProxyGo int64
	newElementsChan     chan T
	enqueueErrorChan    chan error
	// dequeueProxy 协程相关
	numOfDequeueProxyGo         int64
	expiredElements             chan T
	wakeupSignalForDequeueProxy chan struct{}
}

func NewDelayQueue[T Delayable[T]](capacity int, compare ekit.Comparator[T]) (*DelayQueue[T], error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("%w: capacity必须大于0", errInvalidArgument)
	}
	if compare == nil {
		return nil, fmt.Errorf("%w: compare不能为nil", errInvalidArgument)
	}
	d := &DelayQueue[T]{
		mutex:                &sync.RWMutex{},
		q:                    queue.NewPriorityQueue[T](capacity, compare),
		capacity:             capacity,
		compareFuncOfElement: compare,
		// 代理协程生命周期管理相关
		startedProxiesWaitGroup: &sync.WaitGroup{},
		stoppedProxiesWaitGroup: &sync.WaitGroup{},
		quitSignalChan:          make(chan struct{}),
		// enqueueProxy
		newElementsChan:  make(chan T),
		enqueueErrorChan: make(chan error),
		// dequeueProxy
		// expiredElements 必须有缓冲区
		expiredElements:             make(chan T, capacity),
		wakeupSignalForDequeueProxy: make(chan struct{}, 1),
	}
	d.startProxies()
	return d, nil
}

func (d *DelayQueue[T]) startProxies() {
	proxies := 2
	d.stoppedProxiesWaitGroup.Add(proxies)
	d.startedProxiesWaitGroup.Add(proxies)
	go d.enqueueProxy()
	go d.dequeueProxy()
	d.startedProxiesWaitGroup.Wait()
	atomic.AddInt64(&d.numOfEnqueueProxyGo, 1)
	atomic.AddInt64(&d.numOfDequeueProxyGo, 1)
	// log.Println("Proxies Started .....")
}

func (d *DelayQueue[T]) enqueueProxy() {

	defer func() {
		// 吞掉panic，使协程正常退出
		_ = recover()
		d.stoppedProxiesWaitGroup.Done()
		// log.Println("enqueueProxy stop....")
	}()

	d.startedProxiesWaitGroup.Done()
	// log.Println("enqueueProxy start....")

	for {

		select {
		case <-d.quitSignalChan:
			return
		case e := <-d.newElementsChan:
			// log.Println("enqueueProxy, get element ", e)
			if !d.isValidElement(e) {
				// log.Println("enqueueProxy, send err == invalid t ... ")
				d.enqueueErrorChan <- fmt.Errorf("%w: 元素t非法", errInvalidArgument)
				// log.Println("enqueueProxy, blocking... ")
				continue
			}

			d.mutex.Lock()
			// 队列已满
			// 注意：不能将其移动到锁之外
			isFull := d.q.Len()+len(d.expiredElements) == d.capacity
			if isFull {
				d.mutex.Unlock()
				// log.Println("enqueueProxy, send err == Full ... ")
				d.enqueueErrorChan <- queue.ErrOutOfCapacity
				// log.Println("enqueueProxy, blocking... ")
				continue
			}
			// todo: 优化点：为 d.newElementsChan 设置缓冲区，拿到一次锁Enqueue5-10个，过多会饿死 dequeueProxy
			//       需要注意容量判断问题，详见上方isFull
			_ = d.q.Enqueue(e)

			// 写锁保护中，刚入队成功，一定能拿到
			head, _ := d.q.Peek()

			// 新入队元素e具有相等或更高优先级，等于0为了兼容队列为空的情况，并且没有未接收的信号，才考虑发送唤醒信号
			headOfQueueHasChanged := d.compareFuncOfElement(e, head) <= 0
			thereIsNoUnreceivedWakeupSignal := len(d.wakeupSignalForDequeueProxy) == 0
			if headOfQueueHasChanged && thereIsNoUnreceivedWakeupSignal {
				d.wakeupSignalForDequeueProxy <- struct{}{}
				// log.Println("enqueueProxy, notify dequeueProxy ... ")
			}

			d.mutex.Unlock()

			// 通知 Enqueue 协程入队成功
			d.enqueueErrorChan <- (error)(nil)
			// log.Println("enqueueProxy, send err == nil , element enqueued ....", e, "len = ", d.Len())
		}
	}
}

func (d *DelayQueue[T]) isValidElement(elem T) (ok bool) {
	defer func() {
		ok = recover() == nil
	}()
	// elem 是 (*xxx)(nil)
	// 或者 Deadline方法内直接panic
	_ = elem.Deadline()
	return true
}

func (d *DelayQueue[T]) dequeueProxy() {

	defer func() {
		// 吞掉panic，使协程正常退出
		_ = recover()
		d.stoppedProxiesWaitGroup.Done()
		// log.Println("dequeueProxy stop....")
	}()

	d.startedProxiesWaitGroup.Done()
	// log.Println("dequeueProxy start....")

	defaultBlockingDuration := time.Hour
	ticker := time.NewTicker(defaultBlockingDuration)
	defer ticker.Stop()
	var remainingBlockingDuration time.Duration

	for {
		// log.Println("dequeueProxy, peek before....")
		d.mutex.RLock()
		head, err := d.q.Peek()
		d.mutex.RUnlock()
		// log.Println("dequeueProxy, peek after....")

		// 队列为空
		if err != nil {
			// log.Println("dequeueProxy, blocking before, queue empty....")
			ticker.Reset(defaultBlockingDuration)
			goto blocking
		}

		remainingBlockingDuration = time.Duration(head.Deadline().Unix() - time.Now().Unix())
		if remainingBlockingDuration > 0 {
			// 数据未过期
			// log.Println("dequeueProxy, blocking before, waiting duration ....", head, remainingBlockingDuration, time.Now())
			ticker.Reset(remainingBlockingDuration)
			goto blocking
		}

		// 数据已过期
		// log.Println("dequeueProxy, send before....")
		d.mutex.Lock()
		head, _ = d.q.Dequeue()
		d.expiredElements <- head
		// log.Println("dequeueProxy, element dequeued .... ", head, d.q.Len())
		d.mutex.Unlock()
		// log.Println("dequeueProxy, send after....")
		continue

	blocking:
		// log.Println("dequeueProxy, blocking....")
		select {
		case <-d.quitSignalChan:
			return
		case <-ticker.C:
			// log.Println("dequeueProxy, unblocked by Ticker.....")
		case <-d.wakeupSignalForDequeueProxy:
			// log.Println("dequeueProxy, unblocked by Signal.....")
		}
	}
}

func (d *DelayQueue[T]) Enqueue(ctx context.Context, t T) error {

	if d.isClosed() {
		return fmt.Errorf("%w", errQueueHasBeenClosed)
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	// log.Println("Enqueue, Waiting for adding element....")
	select {
	case <-ctx.Done():
		return ctx.Err()
	case d.newElementsChan <- t:
		return <-d.enqueueErrorChan
	}
}

func (d *DelayQueue[T]) Dequeue(ctx context.Context) (T, error) {

	var zeroValue T

	if d.isClosed() {
		return zeroValue, fmt.Errorf("%w", errQueueHasBeenClosed)
	}

	if ctx.Err() != nil {
		return zeroValue, ctx.Err()
	}

	// log.Println("Dequeue, Waiting for element....")

	select {
	case <-ctx.Done():
		return zeroValue, ctx.Err()
	case elem := <-d.expiredElements:
		// log.Println("Dequeue ...", elem)
		return elem, nil
	}
}

func (d *DelayQueue[T]) isClosed() bool {
	return atomic.LoadInt64(&d.numOfEnqueueProxyGo) == 0 && atomic.LoadInt64(&d.numOfDequeueProxyGo) == 0
}

func (d *DelayQueue[T]) Len() int {
	d.mutex.RLock()
	// 一部分过期数据会缓存在 expiredElements 中
	// 但并未被Dequeue调用协程取走，所以逻辑上还是要将缓存数据算在内的。
	n := d.q.Len() + len(d.expiredElements)
	d.mutex.RUnlock()
	return n
}

func (d *DelayQueue[T]) Close() {
	if atomic.CompareAndSwapInt64(&d.numOfEnqueueProxyGo, 1, 0) &&
		atomic.CompareAndSwapInt64(&d.numOfDequeueProxyGo, 1, 0) {
		close(d.quitSignalChan)
		d.stoppedProxiesWaitGroup.Wait()
		// log.Println("Proxies Stopped .....")
	}
}
