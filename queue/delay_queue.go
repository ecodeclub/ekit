package queue

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gotomicro/ekit"
	"github.com/gotomicro/ekit/internal/queue"
)

type Delayable[T any] interface {
	Delay() time.Duration
}

type DelayQueue[T Delayable[T]] struct {
	q     *queue.PriorityQueue[T]
	mutex *sync.RWMutex

	compareFunc ekit.Comparator[T]

	headOfQueueHasChangedSignal chan struct{}
	SignalSenders               int64

	SignalReceivers int64
}

func NewDelayQueue[T Delayable[T]](compare ekit.Comparator[T]) *DelayQueue[T] {
	d := &DelayQueue[T]{
		q:                           queue.NewPriorityQueue[T](0, compare),
		mutex:                       &sync.RWMutex{},
		headOfQueueHasChangedSignal: make(chan struct{}),
	}
	return d
}

func (d *DelayQueue[T]) Enqueue(ctx context.Context, t T) error {

	if ctx.Err() != nil {
		return ctx.Err()
	}

	d.mutex.Lock()

	// 排队超时
	if ctx.Err() != nil {
		d.mutex.Unlock()
		return ctx.Err()
	}

	// 无界队列，不会返回err
	_ = d.q.Enqueue(t)

	// 写锁保护，刚刚入队，一定能拿到
	head, _ := d.q.Peek()

	// 元素t具有相等或更高优先级，等于0为了兼容队列为空的情况
	if d.compareFunc(t, head) <= 0 {
		atomic.AddInt64(&d.SignalSenders, 1)
		defer atomic.AddInt64(&d.SignalSenders, -1)
	}

	d.mutex.Unlock()

	trySendSignalCounter := 0
	// 有并发协程调用Dequeue才需要发送信号
	// 发送原则是：尽最大努力发送信号，不要永久阻塞自己，因为当Dequeue上的并发协程数大于队列中元素数还是有会阻塞的
	// 同一时间段内排队的协程中，只有最后一个协程才有资格发送信号，记作：g10
	// 下一时间段内排队的协程中，g20有资格发送信号，但因g10的存在而无法发送信号进而直接退出
	if atomic.LoadInt64(&d.SignalReceivers) > 0 && atomic.LoadInt64(&d.SignalSenders) == 1 {
		// 当g10执行到这里，d.wakeUpSignal缓冲区为0，恰好阻塞在Dequeue中select上的协程走case <-d.wakeUpSignal分支，g10退出
		// 当g30执行到这里，d.wakeUpSignal缓冲区为0，恰好阻塞在Dequeue中select上的协程走其他case，此时
		//   如果设置了超时，g30将在超时时间内，尽最大努力发送信号
		//   如果未设置超时，g30将在尝试一定次数之后退出。
		select {
		case <-ctx.Done():
			// g30因超时离开
			return ctx.Err()
		case d.headOfQueueHasChangedSignal <- struct{}{}:
			// g10或g30因发送信号成功而离开
			return nil
		default:
			trySendSignalCounter++
			if trySendSignalCounter == 10 {
				// g30因达到最大尝试次数而离开
				return nil
			}
		}
	}
	return nil
}

type Cond struct {
	sync.Cond
}

func (d *DelayQueue[T]) Dequeue(ctx context.Context) (T, error) {

	atomic.AddInt64(&d.SignalReceivers, 1)
	defer atomic.AddInt64(&d.SignalReceivers, -1)

	var zeroValue T

	// 过期超时
	if ctx.Err() != nil {
		return zeroValue, ctx.Err()
	}

	ticker := time.NewTicker(0)
	ticker.Stop()

	for {

		d.mutex.Lock()

		// 排队超时
		if ctx.Err() != nil {
			d.mutex.Unlock()
			return zeroValue, ctx.Err()
		}

		head, err := d.q.Peek()

		d.mutex.Unlock()

		if err != nil {
			select {
			case <-ctx.Done():
				return zeroValue, ctx.Err()
			case <-d.headOfQueueHasChangedSignal:
			}
			continue
		}

		// todo: head是指针还是其他类型，如果head是指针，那么head可能已经被出队
		duration := head.Delay() - time.Duration(time.Now().UnixNano())
		if duration <= 0 {
			continue
		}

		ticker.Reset(duration)

		select {
		case <-d.headOfQueueHasChangedSignal:
			// 最坏情况，多个协程阻塞，从大到小一次入队，100s，50s, 25s, 12, 6, 3, 1
			// 只唤醒一个，去队头重新取元素
		case <-ctx.Done():
			return zeroValue, ctx.Err()
		case <-ticker.C:
			d.mutex.Lock()
			// 再次检查队头
			// 1）此时恰好有更高优先级的元素入队，对延迟队列逻辑语义无影响，可以直接出队。
			// 2）原队头被前一个协程拿走，而新队头与原队头之间较大时间差，不能直接出队。
			//    当前协程需要再次阻塞，可能会出现饿死的情况。
			head, err := d.q.Peek()
			if err == nil {
				// 验证队头元素确实过期
				duration := head.Delay() - time.Duration(time.Now().UnixNano())
				if duration <= 0 {
					// 一定成功
					t, err := d.q.Dequeue()
					d.mutex.Unlock()
					return t, err
				}
				// todo：优化点
				//   协程经历过X次阻塞后，且duration大于0且小于阀值Y，在不释放写锁的情况下直接阻塞让其拿到队头
				//   以解决饿死的问题
			}
			d.mutex.Unlock()
			// 程序走到这，表示原队头被拿走，需要再次去队头
		}
	}
}
