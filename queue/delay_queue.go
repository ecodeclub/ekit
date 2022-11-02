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
	stateOfProxyGoroutinesCreated int64 = 1
	stateOfProxyGoroutinesRunning int64 = 2
	stateOfProxyGoroutinesStopped int64 = 3

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
	qCache               *queue.PriorityQueue[T]
	capacity             int
	compareFuncOfElement ekit.Comparator[T]

	// Proxy 协程生命周期管理相关
	startedProxiesWaitGroup *sync.WaitGroup
	stoppedProxiesWaitGroup *sync.WaitGroup
	quitSignalChan          chan struct{}
	stateOfProxyGoroutines  int64
	// enqueueProxy 协程相关
	enqueueElemsChan chan T
	enqueueErrorChan chan error
	// workerProxy 协程相关
	wakeupWorkerProxy chan struct{}
	// dequeueProxy 协程相关
	dequeueElemsChan   chan T
	dequeueErrorChan   chan error
	wakeupDequeueProxy chan struct{}
}

func NewDelayQueue[T Delayable[T]](capacity int) (*DelayQueue[T], error) {
	if capacity < 0 {
		return nil, fmt.Errorf("%w: capacity必须大于等于0", errInvalidArgument)
	}

	compare := func(t1 T, t2 T) int {
		t1Unix := t1.Deadline().Unix()
		t2Unix := t2.Deadline().Unix()
		if t1Unix < t2Unix {
			return -1
		} else if t1Unix == t2Unix {
			return 0
		} else {
			return 1
		}
	}

	d := &DelayQueue[T]{
		mutex:                &sync.RWMutex{},
		q:                    queue.NewPriorityQueue[T](capacity, compare),
		qCache:               queue.NewPriorityQueue[T](capacity, compare),
		capacity:             capacity,
		compareFuncOfElement: compare,
		// 代理协程生命周期管理相关
		startedProxiesWaitGroup: &sync.WaitGroup{},
		stoppedProxiesWaitGroup: &sync.WaitGroup{},
		quitSignalChan:          make(chan struct{}),
		stateOfProxyGoroutines:  stateOfProxyGoroutinesCreated,
		// enqueueProxy
		enqueueElemsChan: make(chan T),
		enqueueErrorChan: make(chan error),
		// workerProxy
		wakeupWorkerProxy: make(chan struct{}, 1),
		// dequeueProxy
		dequeueElemsChan:   make(chan T),
		dequeueErrorChan:   make(chan error),
		wakeupDequeueProxy: make(chan struct{}, 1),
	}
	d.startProxies()
	return d, nil
}

func (d *DelayQueue[T]) startProxies() {
	proxies := 3
	d.stoppedProxiesWaitGroup.Add(proxies)
	d.startedProxiesWaitGroup.Add(proxies)
	go d.enqueueProxy()
	go d.workerProxy()
	go d.dequeueProxy()
	d.startedProxiesWaitGroup.Wait()
	atomic.StoreInt64(&d.stateOfProxyGoroutines, stateOfProxyGoroutinesRunning)
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
		case e := <-d.enqueueElemsChan:
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
			isFull := d.capacity > 0 && d.q.Len()+d.qCache.Len() == d.capacity
			if isFull {
				d.mutex.Unlock()
				// log.Println("enqueueProxy, send err == Full ... ")
				d.enqueueErrorChan <- queue.ErrOutOfCapacity
				// log.Println("enqueueProxy, blocking... ")
				continue
			}

			// todo: 优化点：为 d.enqueueElemsChan 设置缓冲区，拿到一次锁Enqueue5-10个，过多会饿死 workerProxy
			//       需要注意容量判断问题，详见上方isFull
			_ = d.q.Enqueue(e)

			// 写锁保护中，刚入队成功，一定能拿到
			head, _ := d.q.Peek()

			// 新入队元素e具有相等或更高优先级，等于0为了兼容队列为空的情况，并且没有未接收的信号，才考虑发送唤醒信号
			headOfQueueHasChanged := d.compareFuncOfElement(e, head) <= 0
			thereIsNoUnreceivedWakeupSignal := len(d.wakeupWorkerProxy) == 0
			if headOfQueueHasChanged && thereIsNoUnreceivedWakeupSignal {
				d.wakeupWorkerProxy <- struct{}{}
				// log.Println("enqueueProxy, notify workerProxy ... ")
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

func (d *DelayQueue[T]) workerProxy() {

	defer func() {
		// 吞掉panic，使协程正常退出
		_ = recover()
		d.stoppedProxiesWaitGroup.Done()
		// log.Println("workerProxy stop....")
	}()

	d.startedProxiesWaitGroup.Done()
	// log.Println("workerProxy start....")

	defaultBlockingDuration := time.Hour
	ticker := time.NewTicker(defaultBlockingDuration)
	defer ticker.Stop()
	var remainingBlockingDuration time.Duration

	for {
		// log.Println("workerProxy, peek before....")
		d.mutex.RLock()
		head, err := d.q.Peek()
		d.mutex.RUnlock()
		// log.Println("workerProxy, peek after....")

		// 队列为空
		if err != nil {
			// log.Println("workerProxy, blocking before, queue empty....")
			ticker.Reset(defaultBlockingDuration)
			goto blocking
		}

		remainingBlockingDuration = time.Duration(head.Deadline().Unix() - time.Now().Unix())
		if remainingBlockingDuration > 0 {
			// 数据未过期
			// log.Println("workerProxy, blocking before, waiting duration ....", head, remainingBlockingDuration, time.Now())
			ticker.Reset(remainingBlockingDuration)
			goto blocking
		}

		// 数据已过期
		// log.Println("workerProxy, cache before....")
		d.mutex.Lock()
		head, _ = d.q.Dequeue()
		// 将过期数据放入缓存中，保持顺序
		_ = d.qCache.Enqueue(head)
		// log.Println("workerProxy, element dequeued from q enqueue into qCache .... ", head, d.q.Len(), d.qCache.Len())
		if d.qCache.Len() == 1 && len(d.wakeupDequeueProxy) == 0 {
			d.wakeupDequeueProxy <- struct{}{}
		}
		d.mutex.Unlock()
		// log.Println("workerProxy, cache after....")
		continue

	blocking:
		// log.Println("workerProxy, blocking....")
		select {
		case <-d.quitSignalChan:
			return
		case <-ticker.C:
			// log.Println("workerProxy, unblocked by Ticker.....")
		case <-d.wakeupWorkerProxy:
			// log.Println("workerProxy, unblocked by Signal.....")
		}
	}
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

	for {
		// log.Println("dequeueProxy, peeked before...")
		d.mutex.RLock()
		elem, err := d.qCache.Peek()
		d.mutex.RUnlock()
		// log.Println("dequeueProxy, peeked after...")
		if err != nil {
			// log.Println("dequeueProxy, blocking, qCache empty ...")
			select {
			case <-d.quitSignalChan:
				return
			case <-d.wakeupDequeueProxy:
				// log.Println("dequeueProxy, unblocked by wakeup ...")
			}
			continue
		}
		// log.Println("dequeueProxy, peeked element from qCache .... ", elem)

		select {
		case <-d.quitSignalChan:
			return
		case d.dequeueElemsChan <- elem:
			d.mutex.Lock()
			// 元素一定存在，且err == nil
			_, err := d.qCache.Dequeue()
			// log.Println("dequeueProxy, element dequeued from qCache .... ", head, d.qCache.Len()+d.q.Len())
			d.mutex.Unlock()

			d.dequeueErrorChan <- err
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
	case d.enqueueElemsChan <- t:
		err := <-d.enqueueErrorChan
		// log.Println("Enqueue, Get response ....", err)
		return err
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
	case elem := <-d.dequeueElemsChan:
		err := <-d.dequeueErrorChan
		// log.Println("Dequeue ...", elem, " len = ", d.Len(), err)
		return elem, err
	}
}

func (d *DelayQueue[T]) isClosed() bool {
	return atomic.LoadInt64(&d.stateOfProxyGoroutines) == stateOfProxyGoroutinesStopped
}

func (d *DelayQueue[T]) Len() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.q.Len() + d.qCache.Len()
}

func (d *DelayQueue[T]) Cap() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.capacity
}

func (d *DelayQueue[T]) Close() {
	if atomic.CompareAndSwapInt64(&d.stateOfProxyGoroutines, stateOfProxyGoroutinesRunning, stateOfProxyGoroutinesStopped) {
		close(d.quitSignalChan)
		d.stoppedProxiesWaitGroup.Wait()
		// log.Println("Proxies Stopped .....")
	}
}
