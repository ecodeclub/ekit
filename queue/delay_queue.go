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

// g1  - \
// g2  - - Enqueue() -- N --> channel（多个通道） -- 1 --> enqueueProxy 协程 --
// g3  - /                                                                   \
//                                                      两个代理协程之间通过互斥锁、通道协作，详见下方说明3
// g4  - \                                                                   /
// g5  - - Dequeue() -- N --> channel（多个通道）-- 1 --> dequeueProxy 协程 --
// g6  - /
//
// 说明：
//    1. g1、g2、g3并发调用Enqueue方法，Enqueue方法内部启动一个代理协程 enqueueProxy，
//       - enqueueProxy 职责
//         - 从调用者协程们接收数据
//         - 与下方 dequeueProxy 协程并发访问底层无锁优先级队列，将收到的数据入队
//         - 将入队结果返回给调用者协程们
//         - 如果有高优先级元素入队，导致队头变更，向 d.wakeupSignalForDequeueProxy 发信号通知下方 dequeueProxy 协程
//         - 监听退出信号，来自最后一个调用者协程发送的退出信号
//       - g1、g2、g3 通过channel与 enqueueProxy 协程通信
//         - 调用协程们通过 d.newElementsChan 向 enqueueProxy 发送数据
//         - 调用协程们通过 d.enqueueErrorChan 从 enqueueProxy 接收错误信息
//         - 第一/最后一个调用协程通过 d.quitSignalForEnqueueProxy 和 d.continueSignalFromEnqueueProxy 启动/关闭 enqueueProxy 协程
//           非并发，g1既是第一个又是最后一个需要负责启动和关闭 enqueueProxy
//           并发下，g1负责启动 enqueueProxy，g3 负责关闭 enqueueProxy
//
//    2. g3、g4、g5并发调用Dequeue方法，Dequeue方法内部启动一个代理协程 dequeueProxy，
//       - dequeueProxy 职责
//         - 获取队头，等待其过期；
// 		- 与上方 enqueueProxy 协程并发访问底层无锁优先级队列，将过期队头出队
//         - 将出队结果返回给调用者协程们
//         - 监听唤醒信号，来自 enqueueProxy 协程，重新检查队头
//         - 监听退出信号，来自最后一个调用者协程发送的退出信号
//       - g3、g4、g5 通过channel与 DequeueProxy 协程通信
//         - 调用者协程们从 d.expiredElements 获取过期元素
//         - dequeueProxy 协程通过 d.wakeupSignalForDequeueProxy 获取通知以重新检查队头
//         - 第一/最后一个调用协程通过 d.quitSignalForDequeueProxy 和 d.continueSignalFromDequeueProxy 启动/关闭 dequeueProxy 协程
//         - Dequeue的逻辑语义是"拿到队头，等待队头超时或自己超时返回; 拿不到队头，阻塞等待直到ctx过期"
//           故Dequeue不返回延迟队列为空的错误，而是让调用者阻塞等待ctx超时；如果调用者未传递具有超时的ctx，导致永久阻塞是他自己的问题
//
//    3. enqueueProxy 与 dequeueProxy 协程之间通过互斥锁、通道来协作
//      - 用 d.mutex 并发操作底层优先级队列 d.q
//      - 用 d.wakeupSignalForDequeueProxy && d.wakeupSignalForEnqueueProxy 在队列状态变化时相互唤醒
//      - 无锁队列 d.q 上只有 enqueueProxy 与 dequeueProxy 两个协程并发访问

var (
	errInvalidArgument = errors.New("ekit: 参数非法")
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

	// Enqueue方法上并发调用者协程计数
	enqueueMutex   *sync.Mutex
	enqueueCallers int64

	// enqueueProxy 协程相关
	numOfEnqueueProxyGo            int64
	newElementsChan                chan T
	enqueueErrorChan               chan error
	quitSignalForEnqueueProxy      chan struct{}
	continueSignalFromEnqueueProxy chan struct{}
	wakeupSignalForEnqueueProxy    chan struct{}

	// Dequeue方法上并发调用者协程计数
	dequeueMutex   *sync.Mutex
	dequeueCallers int64

	// dequeueProxy 协程相关
	numOfDequeueProxyGo            int64
	expiredElements                chan T
	quitSignalForDequeueProxy      chan struct{}
	continueSignalFromDequeueProxy chan struct{}
	wakeupSignalForDequeueProxy    chan struct{}
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

		enqueueMutex: &sync.Mutex{},
		// enqueueProxy
		newElementsChan:                make(chan T),
		enqueueErrorChan:               make(chan error),
		wakeupSignalForEnqueueProxy:    make(chan struct{}, 1),
		quitSignalForEnqueueProxy:      make(chan struct{}, 1),
		continueSignalFromEnqueueProxy: make(chan struct{}, 1),

		dequeueMutex: &sync.Mutex{},
		// dequeueProxy
		// expiredElements 必须有缓冲区
		expiredElements:                make(chan T, capacity),
		wakeupSignalForDequeueProxy:    make(chan struct{}, 1),
		quitSignalForDequeueProxy:      make(chan struct{}, 1),
		continueSignalFromDequeueProxy: make(chan struct{}, 1),
	}
	return d, nil
}

func (d *DelayQueue[T]) Enqueue(ctx context.Context, t T) error {

	if ctx.Err() != nil {
		return ctx.Err()
	}

	// 更新计数与启动/关闭 enqueueProxy 协程必须是原子的
	d.enqueueMutex.Lock()
	d.enqueueCallers++
	// 第一个检测到 enqueueProxy 协程没启动的调用者负责启动 enqueueProxy 协程，调用者进入下方select语言前 enqueueProxy 必须启动
	if atomic.LoadInt64(&d.numOfEnqueueProxyGo) == 0 {
		go d.enqueueProxy()
		// enqueueProxy 协程启动后，通知当前协程继续运行
		<-d.continueSignalFromEnqueueProxy
	}
	d.enqueueMutex.Unlock()

	defer func() {
		d.enqueueMutex.Lock()
		d.enqueueCallers--
		// 最后一个检测到 enqueueProxy 协程存在的调用者，在退出前需要确保 enqueueProxy 协程先于自己退出。
		if d.enqueueCallers == 0 && atomic.LoadInt64(&d.numOfEnqueueProxyGo) == 1 {
			// 以非阻塞方式发送信号，通知 enqueueProxy 协程走退出流程，当前协程阻塞等待
			d.quitSignalForEnqueueProxy <- struct{}{}
			<-d.continueSignalFromEnqueueProxy
		}
		d.enqueueMutex.Unlock()
	}()

	// log.Println("Enqueue, Waiting for adding element....")
	select {
	case <-ctx.Done():
		return ctx.Err()
	case d.newElementsChan <- t:
		return <-d.enqueueErrorChan
	}
}

func (d *DelayQueue[T]) enqueueProxy() {

	defer func() {
		// 吞掉panic
		_ = recover()
		// 必须先更新计数，再发送信号
		atomic.CompareAndSwapInt64(&d.numOfEnqueueProxyGo, 1, 0)
		d.continueSignalFromEnqueueProxy <- struct{}{}
		// log.Println("enqueueProxy stop....")
	}()

	// 必须先更新计数，再发送信号
	atomic.CompareAndSwapInt64(&d.numOfEnqueueProxyGo, 0, 1)
	d.continueSignalFromEnqueueProxy <- struct{}{}

	// log.Println("enqueueProxy start....")

	for {

		select {
		case e := <-d.newElementsChan:
			// log.Println("enqueueProxy, get element ", e)

			d.mutex.Lock()

			// 队列已满
			// dequeueProxy 协程在向 d.expiredElements 发送数据时
			// 即 d.dequeueAndSendExpiredElement() 中实现也用到了写锁
			// 所以这里可以认为 len( d.expiredElements ) 不变
			isFull := d.q.Len()+len(d.expiredElements) == d.capacity
			if isFull {
				d.mutex.Unlock()
				// log.Println("enqueueProxy, send err == Full ... ")
				// 通知当前 Enqueue 协程队列已满
				d.enqueueErrorChan <- queue.ErrOutOfCapacity
				// log.Println("enqueueProxy, blocking... ")
				continue
			}
			// todo: 优化点：为 d.newElementsChan 设置缓冲区，拿到一次锁尽可能将缓冲区中数据全部Enqueue
			//       需要注意容量判断问题，详见上方isFull
			_ = d.q.Enqueue(e)
			// err := d.q.Enqueue(e)
			// if err != nil {
			// 	// 队列已满
			// 	d.mutex.Unlock()
			// 	log.Println("enqueueProxy, send err == Full ... ")
			// 	d.enqueueErrorChan <- queue.ErrOutOfCapacity
			// 	log.Println("enqueueProxy, blocking... ")
			// 	continue
			// }

			// 写锁保护中，刚入队成功，一定能拿到
			head, _ := d.q.Peek()
			// 新入队元素e具有相等或更高优先级，等于0为了兼容队列为空的情况
			headOfQueueHasChanged := d.compareFuncOfElement(e, head) <= 0
			// d.wakeupSignalForDequeueProxy 的消费者只有 dequeueProxy 协程
			// 只有 d.wakeupSignalForDequeueProxy 为空时才需要再次发送信号
			// 如果之前的信号还未被 dequeueProxy 协程消费，未消费的信号也能表示相同意义
			// 即 dequeueProxy 协程拿到信号后重新去检查队头
			thereIsNoUnreceivedWakeupSignal := len(d.wakeupSignalForDequeueProxy) == 0

			if headOfQueueHasChanged && thereIsNoUnreceivedWakeupSignal {
				d.wakeupSignalForDequeueProxy <- struct{}{}
				// log.Println("enqueueProxy, notify dequeueProxy ... ")
			}
			d.mutex.Unlock()

			// 通知 Enqueue 协程入队成功
			d.enqueueErrorChan <- (error)(nil)
			// log.Println("enqueueProxy, send err == nil , element enqueued ....", e, "len = ", d.Len())

		case <-d.quitSignalForEnqueueProxy:
			return
			// case <-d.wakeupSignalForEnqueueProxy:
			// 	// 等待 dequeueProxy 在调用 d.q.Dequeue 后发送信号将自己唤醒
			// 	log.Println("enqueueProxy, wakeup by dequeueProxy ... ")
		}
	}
}

func (d *DelayQueue[T]) Dequeue(ctx context.Context) (T, error) {

	var zeroValue T
	if ctx.Err() != nil {
		return zeroValue, ctx.Err()
	}

	// 更新计数与启动/关闭 dequeueProxy 协程必须是原子的
	d.dequeueMutex.Lock()
	d.dequeueCallers++
	// 第一个检测到 dequeueProxy 协程没有启动的协程进入下方select前，需要将 dequeueProxy 协程启动起来
	if atomic.LoadInt64(&d.numOfDequeueProxyGo) == 0 {
		go d.dequeueProxy()
		// dequeueProxy 协程启动后，唤醒当前协程
		<-d.continueSignalFromDequeueProxy
	}
	d.dequeueMutex.Unlock()

	defer func() {
		d.dequeueMutex.Lock()
		d.dequeueCallers--
		// 最后一个检测到 dequeueProxy 协程存在的协程，在退出前需要确保 dequeueProxy 协程先于自己退出。
		if d.dequeueCallers == 0 && atomic.LoadInt64(&d.numOfDequeueProxyGo) == 1 {
			// 以非阻塞方式发送信号，通知 dequeueProxy 协程走退出流程，当前协程阻塞等待
			d.quitSignalForDequeueProxy <- struct{}{}
			// dequeueProxy 协程退出前，通知当前协程退出
			<-d.continueSignalFromDequeueProxy
		}
		d.dequeueMutex.Unlock()
	}()

	// log.Println("Dequeue, Waiting for element....")

	select {
	case <-ctx.Done():
		return zeroValue, ctx.Err()
	case elem := <-d.expiredElements:
		// log.Println("Dequeue ...", elem)
		return elem, nil
	}
}

func (d *DelayQueue[T]) dequeueProxy() {

	defer func() {
		// 吞掉panic
		_ = recover()
		// 必须先更新计数，再发送信号
		atomic.CompareAndSwapInt64(&d.numOfDequeueProxyGo, 1, 0)
		d.continueSignalFromDequeueProxy <- struct{}{}
		// log.Println("dequeueProxy stop....")
	}()

	// 必须先更新计数，再发送信号
	atomic.CompareAndSwapInt64(&d.numOfDequeueProxyGo, 0, 1)
	d.continueSignalFromDequeueProxy <- struct{}{}
	// log.Println("dequeueProxy start....")

	defaultBlockingDuration := time.Hour
	ticker := time.NewTicker(defaultBlockingDuration)
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
		} else {
			// 数据已过期
			// log.Println("dequeueProxy, send before....")
			d.dequeueAndSendExpiredElement()
			// log.Println("dequeueProxy, send after....")
			// 重新获取队头
			continue
		}

	blocking:
		// log.Println("dequeueProxy, blocking....")
		select {
		case <-d.quitSignalForDequeueProxy:
			return
		case <-ticker.C:
			// 因为只有 dequeueProxy 协程调用d.q.Dequeue()
			// 阻塞醒来后是可以立即调用d.q.Dequeue()，即便与b.q.Enqueue()并发也是正确的。
			// 不存在原队头被其他协程并发调用d.q.Dequeue()出队的情况(前提单 dequeueProxy 实例）
			// 为了便于理解程序，将下方调用代码注释，不注释下方调用的逻辑分析见方法内部注释
			// d.dequeueAndSendExpiredElement()
			// log.Println("dequeueProxy, unblocked by Ticker.....")
		case <-d.wakeupSignalForDequeueProxy:
			// 队头更新，再次检查
			// log.Println("dequeueProxy, unblocked by Signal.....")
		}
	}
}

func (d *DelayQueue[T]) dequeueAndSendExpiredElement() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	// 前置条件：单 dequeueProxy 协程实例

	// 如果不注释掉 dequeueProxy 中case <-ticker.C下的本方法调用，在单 dequeueProxy 实例的前提下，程序执行到这里有两种情况：
	// Case 1. 队头元素X过期，直接出队元素X
	// Case 2. 协程按原队头X的 remainingBlockingDuration 阻塞中：
	// 	        A. 阻塞到期后，执行d.q.Dequeue()删除原队头X
	// 	        B. 收到新元素Y入队的信号，重新去队头获取阻塞时间duration
	//             1) 如果已过期duration <= 0，直接将新队头出队（即重复[Case 1]操作）此时出队的是新队头Y
	//             2) 如果还需阻塞，重复[Case 2.A]操作，此时出队的也是新队头Y
	//             3) 如果有新元素入队，此时重复[Case 2.B]操作
	// 通过以上分析可知，无论d.q.Dequeue()出队的是原队头X还是新队头Y/Z/M等(多次插队），因插队元素相比于原队头具有更高或相等的优先级
	// 那么即便按照原队头的remainingBlockingDuration进行阻塞且被唤醒后直接d.q.Dequeue()删除的是新队头也是安全。
	// 因为新队头的截止日期应该早于或等原队头的截止日期。

	// 如果注释掉 dequeueProxy 中case <-ticker.C下的本方法调用，在单 dequeueProxy 实例的前提下，进入这里只有一种情况：remainingBlockingDuration <= 0

	// 单 dequeueProxy 协程下，一定没问题
	// 多 dequeueProxy 协程下，要检查当前队头与阻塞前获取的队头是一样的，不一样要
	expired, _ := d.q.Dequeue()

	// 非阻塞
	// 即便此时 dequeueProxy 协程收到退出信号,因 d.expiredElements 有缓冲区且与 d.q 的容量相同
	// 因此过期的数据会缓存在 d.expiredElements 中后续调用Dequeue的协程可以直接拿到
	// 注意: d.expiredElements 一定要有缓冲，
	// 否则唯一的Dequeue调用者协程因ctx超时走退出流程时，
	// Dequeue 调用者协程等待在 d.quitSignalForDequeueProxy 上而 dequeueProxy 协程等待在 d.expiredElements <- expired 上从而形成两者相互等待
	d.expiredElements <- expired
	// log.Println("dequeueProxy, element dequeued .... ", expired, d.q.Len())

	// 之前容量为满，enqueueProxy 可能被阻塞
	// enqueueProxyMayBeBlocked := d.q.Len()+len(d.expiredElements) == d.capacity
	// thereIsNoUnreceivedWakeupSignal := len(d.wakeupSignalForEnqueueProxy) == 0
	// if enqueueProxyMayBeBlocked && thereIsNoUnreceivedWakeupSignal {
	// 	d.wakeupSignalForEnqueueProxy <- struct{}{}
	// }
}

func (d *DelayQueue[T]) Len() int {
	d.mutex.RLock()
	// 一部分过期数据会缓存在 expiredElements 中
	// 但并未被Dequeue调用协程取走，所以逻辑上还是要将缓存数据算在内的。
	n := d.q.Len() + len(d.expiredElements)
	d.mutex.RUnlock()
	return n
}
