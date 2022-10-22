package queue_test

import (
	"context"
	"testing"
	"time"

	"github.com/gotomicro/ekit/queue"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestDelayQueue_Enqueue(t *testing.T) {

	q := newDelayQueue[*Int]()

	t.Run("超时控制", func(t *testing.T) {

		t.Run("直接超时", func(t *testing.T) {

			cancelCtx, cancelFunc := context.WithCancel(context.Background())
			cancelFunc()
			assert.ErrorIs(t, q.Enqueue(cancelCtx, newInt(1, time.Second)), context.Canceled)

			timeoutCtx, timeoutCancelFunc := context.WithTimeout(context.Background(), time.Nanosecond)
			assert.Equal(t, q.Enqueue(timeoutCtx, newInt(2, time.Second)), context.DeadlineExceeded)
			timeoutCancelFunc()

			deadlineCtx, deadlineCancelFunc := context.WithDeadline(context.Background(), time.Now().Add(1))
			assert.Equal(t, q.Enqueue(deadlineCtx, newInt(3, time.Second)), context.DeadlineExceeded)
			deadlineCancelFunc()
		})

		t.Run("排队超时", func(t *testing.T) {

			var eg errgroup.Group

			n := 50
			waitChan := make(chan struct{})
			for i := 0; i < n; i++ {
				i := i
				eg.Go(func() error {
					<-waitChan
					timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), 100*time.Nanosecond)
					defer cancelFunc()
					return q.Enqueue(timeoutCtx, newInt(i, time.Millisecond))
				})
			}

			// waitChan <- struct{}{}
			close(waitChan)
			assert.Equal(t, context.DeadlineExceeded, eg.Wait())

		})

		t.Run("等待超时", func(t *testing.T) {

		})
	})

	t.Run("单协程", func(t *testing.T) {

		t.Run("出队", func(t *testing.T) {

			t.Run("直接超时", func(t *testing.T) {

			})

			t.Run("排队超时", func(t *testing.T) {

			})

			t.Run("等待超时", func(t *testing.T) {

			})
		})

		t.Run("延迟", func(t *testing.T) {

		})

	})

}

// Dequeue
// Case 1:
//  queue = {100},
//   A. g1 Dequeue，阻塞（在锁中，去阻塞的路上，阻塞中）
//      1）超时退出
//      2）获取数据成功；
//      3）永久阻塞 —— 不需要处理，使用者不设置超时导致
//   B. g1，g2 ...gX Dequeue，阻塞（在锁中，去阻塞的路上，阻塞中）
//      1）全部超时退出
//      2）部分获取数据，部分超时退出
//      3）部分获取数据成功，其他永久阻塞 —— 不需要处理，使用者不设置超时导致
//   C. g1 Dequeue, 阻塞中，
//      1） g2 Enqueue(120), queue = {100, 120}，g1 不受影响阻塞不中断
//      2） g2 Enqueue(1), queue = {1, 100}, g1被唤醒，重新获取队头，再次阻塞指定时间；
//          这里"1"的定义应该是从Enqueue->Dequeue出元素的最大时间，与之前阻塞多久没关系。
//   D. g1，g2，g3 Dequeue，阻塞中
//      1）g4 Enqueue(120),queue = {100, 120}, g1，g2，g3继续阻塞不受影响（重复多次）；
//      2) g4 Enqueue(10), queue = {10, 120},g1/g2/g3被唤醒，再次阻塞指定时间（重复多次）
//         a. 如果唤醒一个，因为ctx检查点的存在，可能导致刚被唤醒的协程直接退出，其他两个协程仍然在阻塞中
//         b. 如果唤醒全部，同上，仍有可能全部协程因ctx检查点退出，也可能导致大量协程（1k/1w）被唤醒但做的是无用功。
//                唤醒全部，只能增加调度开销，但更新了所有协程的等待时间，增加了获取最新时机
//                queue = {100, 120}, enqueue（1），唤醒全部协程，恰巧再次select的ticker和ctx.Done时因超时退出，
//                这种情况，与queue={1}，一直没人调Dequeue情况一致。
//         c. 只要进入Dequeue就表明你有获取数据的意愿，至少上你尝试获取数据一次，去掉for循环中检查ctx检查点，
//            只让其在和ticker/信号队列的select中退出。问题就是超时时间不准确稍稍之后，或者刚拿到最新队头信息，超时了。
//         d. 不公平问题，因调度的随机性（chan先检查对面有没有人再放入缓冲去正所谓来的早，不如来的巧)，
//            先Dequeue的协程不一定能获取数据，可能稍稍滞后的超时退出，也可能永久阻塞（调用者的使用问题，提供超时机制你不用，不能怪设计者）
// Case 2:
//  queue = {10, 5},
//   A. 与Case1.A相同
//   B. g1, g2, ...gX Dequeue，阻塞（在锁中，去阻塞的路上，阻塞中）
//      1）都超时退出
//      2）2个获取数据成功，其他超时退出
//      3）2个获取数据成功，其他永久阻塞 —— 不需要处理，使用者不设置超时导致
// Case 3:
//  queue = {10, 5, 2}
//   A. 退化为Case1.A 或者 Case2.A
//   B. 退化为Case1.B 或者 Case2.B
// Case 4:
//  queue = {},
//   A. g1 Dequeue, 阻塞（在锁中，去阻塞的路上，阻塞中）；g2，Enqueue(X)成功，通知g1，g1获取数据X成功！
//   B. g1,g2 Dequeue, 阻塞（在锁中，去阻塞的路上，阻塞中）；g3，Enqueue(Y)成功，通知g1或g2，
//      1）都因超时而退出
//      2）一个获取数据Y成功，另一个超时退出
//      3) 一个获取数据Y成功，另一个阻塞直到下一次Enqueue(Z),退化为Case4.A情况
//         要考虑队列内部状态，比如通知信号队列中是空/满等,最好使用使用固定长度的队列，这样信号队列的容量可以为capacity
//  C. g1, g2 Dequeue，阻塞（在锁中，去阻塞的路上，阻塞中）；g3 Dequeue 阻塞； g4，Enqueue(10, 12, 15)成功发送"1次"通知给g1/g2/g3，
//     0) bug:有一个协程被唤醒，去Peek，在select的ticker和ctx.Done()阻塞时，可能因ctx.Done()退出，而导致无法将队头元素按时取出
//        解决方法：采取异步通知go func(){}，唤醒全部等待在Dequeue上的协程数（考虑信号队列的capacity来确定发送信号的数量）。
//     1) 全都超时退出
//        拿到信号了，但因再次阻塞时超时而退出，正如次才要唤醒全部等待者，一来，更新所有等待时间；二来，增大获取数据的几率；
//     2) 全部获取数据成功
//     3）部分获取数据成功，部分阻塞 queue = {15} 因下一个Enqueue(8)的通知，queue = {8, 15}
//        退化称为Case2.B —— 可能触发Bug，只唤醒一个，丢失信号，需要唤醒多个或这全部，以保证协程更新阻塞时间，阻塞时间为0时即可拿到数据出队；
//   C.  g1, g2 Dequeue，阻塞（在锁中，去阻塞的路上，阻塞中）；g3 Dequeue ，阻塞 g4，Enqueue(10, 2)成功，发送"两次"（多次）通知给g1/g2/g3，
//       1) 全部因超时而退出
//        拿到信号了，但因再次阻塞时超时而退出，正如次才要唤醒全部等待者，一来，更新所有等待时间；二来，增大获取数据的几率；
//       2) 全部成功获取数据
//       3）部分获取数据成功，部分超时退出
//       4）部分成功获取数据，部分阻塞
//          a. 不再有Enqueue并发调用，退化为Case2
//          b. 还是会有Enqueue并发调用，退化为Case1.D

// 本质上来说，

type Int struct {
	id       int
	deadline time.Duration
}

func newInt(id int, expire time.Duration) *Int {
	return &Int{id: id, deadline: time.Duration(time.Now().Add(expire).Unix())}
}

func (i *Int) Delay() time.Duration {
	return i.deadline
}

func newDelayQueue[T queue.Delayable[T]]() *queue.DelayQueue[T] {
	return queue.NewDelayQueue[T](func(t1 T, t2 T) int {
		if int64(t1.Delay()) < int64(t2.Delay()) {
			return -1
		} else if int64(t1.Delay()) == int64(t2.Delay()) {
			return 0
		} else {
			return 1
		}
	})
}

func TestUnixTime(t *testing.T) {
	delay := 10 * time.Second
	now := time.Now()
	deadline := now.Add(delay)
	assert.Equal(t, 10*time.Second, time.Duration(deadline.UnixNano()-now.UnixNano()))
}
