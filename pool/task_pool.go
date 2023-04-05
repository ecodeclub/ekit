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

package pool

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ecodeclub/ekit/bean/option"
)

var (
	stateCreated int32 = 1
	stateRunning int32 = 2
	stateClosing int32 = 3
	stateStopped int32 = 4
	stateLocked  int32 = 5

	errTaskPoolIsNotRunning = errors.New("ekit: TaskPool未运行")
	errTaskPoolIsClosing    = errors.New("ekit：TaskPool关闭中")
	errTaskPoolIsStopped    = errors.New("ekit: TaskPool已停止")
	errTaskPoolIsStarted    = errors.New("ekit：TaskPool已运行")
	errTaskIsInvalid        = errors.New("ekit: Task非法")
	errTaskRunningPanic     = errors.New("ekit: Task运行时异常")

	errInvalidArgument = errors.New("ekit: 参数非法")

	_            TaskPool = &OnDemandBlockTaskPool{}
	panicBuffLen          = 2048

	defaultMaxIdleTime = 10 * time.Second
)

// TaskFunc 一个可执行的任务
type TaskFunc func(ctx context.Context) error

// Run 执行任务
// 超时控制取决于衍生出 TaskFunc 的方法
func (t TaskFunc) Run(ctx context.Context) error { return t(ctx) }

// taskWrapper 是Task的装饰器
type taskWrapper struct {
	t Task
}

func (tw *taskWrapper) Run(ctx context.Context) (err error) {
	defer func() {
		// 处理 panic
		if r := recover(); r != nil {
			buf := make([]byte, panicBuffLen)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("%w：%s", errTaskRunningPanic, fmt.Sprintf("[PANIC]:\t%+v\n%s\n", r, buf))
		}
	}()
	return tw.t.Run(ctx)
}

type group struct {
	mp map[int]int
	n  int32
	mu sync.RWMutex
}

func (g *group) isIn(id int) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	_, ok := g.mp[id]
	return ok
}

func (g *group) add(id int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.mp[id]; !ok {
		g.mp[id] = 1
		g.n++
	}
}

func (g *group) delete(id int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.mp[id]; ok {
		g.n--
	}
	delete(g.mp, id)
}

func (g *group) size() int32 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.n
}

// OnDemandBlockTaskPool 按需创建goroutine的并发阻塞的任务池
type OnDemandBlockTaskPool struct {
	// TaskPool内部状态
	state int32

	queue             chan Task
	numGoRunningTasks int32

	totalGo int32
	mutex   sync.RWMutex

	// 初始协程数
	initGo int32
	// 核心协程数
	coreGo int32
	// 最大协程数
	maxGo int32
	// 超时组
	timeoutGroup *group
	// 最大空闲时间
	maxIdleTime time.Duration
	// 队列积压率
	queueBacklogRate float64
	shutdownOnce     sync.Once

	// 协程id方便调试程序
	id int32

	// 外部信号
	shutdownDone chan struct{}
	// 内部中断信号
	shutdownNowCtx    context.Context
	shutdownNowCancel context.CancelFunc
}

// NewOnDemandBlockTaskPool 创建一个新的 OnDemandBlockTaskPool
// initGo 是初始协程数
// queueSize 是队列大小，即最多有多少个任务在等待调度
// 使用相应的Option选项可以动态扩展协程数
func NewOnDemandBlockTaskPool(initGo int, queueSize int, opts ...option.Option[OnDemandBlockTaskPool]) (*OnDemandBlockTaskPool, error) {
	if initGo < 1 {
		return nil, fmt.Errorf("%w：initGo应该大于0", errInvalidArgument)
	}
	if queueSize < 0 {
		return nil, fmt.Errorf("%w：queueSize应该大于等于0", errInvalidArgument)
	}
	b := &OnDemandBlockTaskPool{
		queue:        make(chan Task, queueSize),
		shutdownDone: make(chan struct{}, 1),
		initGo:       int32(initGo),
		coreGo:       int32(initGo),
		maxGo:        int32(initGo),
		maxIdleTime:  defaultMaxIdleTime,
	}

	b.shutdownNowCtx, b.shutdownNowCancel = context.WithCancel(context.Background())
	atomic.StoreInt32(&b.state, stateCreated)

	option.Apply(b, opts...)

	if b.coreGo != b.initGo && b.maxGo == b.initGo {
		b.maxGo = b.coreGo
	} else if b.coreGo == b.initGo && b.maxGo != b.initGo {
		b.coreGo = b.maxGo
	}
	if !(b.initGo <= b.coreGo && b.coreGo <= b.maxGo) {
		return nil, fmt.Errorf("%w : 需要满足initGo <= coreGo <= maxGo条件", errInvalidArgument)
	}

	b.timeoutGroup = &group{mp: make(map[int]int)}

	if b.queueBacklogRate < float64(0) || float64(1) < b.queueBacklogRate {
		return nil, fmt.Errorf("%w ：queueBacklogRate合法范围为[0,1.0]", errInvalidArgument)
	}
	return b, nil
}

func WithQueueBacklogRate(rate float64) option.Option[OnDemandBlockTaskPool] {
	return func(pool *OnDemandBlockTaskPool) {
		pool.queueBacklogRate = rate
	}
}

func WithCoreGo(n int32) option.Option[OnDemandBlockTaskPool] {
	return func(pool *OnDemandBlockTaskPool) {
		pool.coreGo = n
	}
}

func WithMaxGo(n int32) option.Option[OnDemandBlockTaskPool] {
	return func(pool *OnDemandBlockTaskPool) {
		pool.maxGo = n
	}
}

func WithMaxIdleTime(d time.Duration) option.Option[OnDemandBlockTaskPool] {
	return func(pool *OnDemandBlockTaskPool) {
		pool.maxIdleTime = d
	}
}

// Submit 提交一个任务
// 如果此时队列已满，那么将会阻塞调用者。
// 如果因为 ctx 的原因返回，那么将会返回 ctx.Err()
// 在调用 Start 前后都可以调用 Submit
func (b *OnDemandBlockTaskPool) Submit(ctx context.Context, task Task) error {
	if task == nil {
		return fmt.Errorf("%w", errTaskIsInvalid)
	}
	// todo: 用户未设置超时，可以考虑内部给个超时提交
	for {

		if atomic.LoadInt32(&b.state) == stateClosing {
			return fmt.Errorf("%w", errTaskPoolIsClosing)
		}

		if atomic.LoadInt32(&b.state) == stateStopped {
			return fmt.Errorf("%w", errTaskPoolIsStopped)
		}

		task = &taskWrapper{t: task}

		ok, err := b.trySubmit(ctx, task, stateCreated)
		if ok || err != nil {
			return err
		}

		ok, err = b.trySubmit(ctx, task, stateRunning)
		if ok || err != nil {
			return err
		}
	}
}

func (b *OnDemandBlockTaskPool) trySubmit(ctx context.Context, task Task, state int32) (bool, error) {
	// 进入临界区
	if atomic.CompareAndSwapInt32(&b.state, state, stateLocked) {
		defer atomic.CompareAndSwapInt32(&b.state, stateLocked, state)

		// 此处b.queue <- task不会因为b.queue被关闭而panic
		// 代码执行到trySubmit时TaskPool处于lock状态
		// 要关闭b.queue需要TaskPool处于RUNNING状态，Shutdown/ShutdownNow才能成功
		select {
		case <-ctx.Done():
			return false, fmt.Errorf("%w", ctx.Err())
		case b.queue <- task:
			if state == stateRunning && b.allowToCreateGoroutine() {
				b.increaseTotalGo(1)
				id := int(atomic.AddInt32(&b.id, 1))
				go b.goroutine(id)
				// log.Println("create go ", id)
			}
			return true, nil
		default:
			// 不能阻塞在临界区,要给Shutdown和ShutdownNow机会
			return false, nil
		}
	}
	return false, nil
}

func (b *OnDemandBlockTaskPool) allowToCreateGoroutine() bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.totalGo == b.maxGo {
		return false
	}

	// 这个判断可能太苛刻了，经常导致开协程失败，先注释掉
	// allGoShouldBeBusy := atomic.LoadInt32(&b.numGoRunningTasks) == b.totalGo
	// if !allGoShouldBeBusy {
	// 	return false
	// }

	rate := float64(len(b.queue)) / float64(cap(b.queue))
	if rate == 0 || rate < b.queueBacklogRate {
		// log.Println("rate == 0", rate == 0, "rate", rate, " < ", b.queueBacklogRate)
		return false
	}

	// b.totalGo < b.maxGo && rate != 0 && rate >= b.queueBacklogRate
	return true
}

// Start 开始调度任务执行
// Start 之后，调用者可以继续使用 Submit 提交任务
func (b *OnDemandBlockTaskPool) Start() error {

	for {

		if atomic.LoadInt32(&b.state) == stateClosing {
			return fmt.Errorf("%w", errTaskPoolIsClosing)
		}

		if atomic.LoadInt32(&b.state) == stateStopped {
			return fmt.Errorf("%w", errTaskPoolIsStopped)
		}

		if atomic.LoadInt32(&b.state) == stateRunning {
			return fmt.Errorf("%w", errTaskPoolIsStarted)
		}

		if atomic.CompareAndSwapInt32(&b.state, stateCreated, stateLocked) {

			n := b.initGo

			allowGo := b.maxGo - b.initGo
			needGo := int32(len(b.queue)) - b.initGo
			if needGo > 0 {
				if needGo <= allowGo {
					n += needGo
				} else {
					n += allowGo
				}
			}

			b.increaseTotalGo(n)
			for i := int32(0); i < n; i++ {
				go b.goroutine(int(atomic.AddInt32(&b.id, 1)))
			}
			atomic.CompareAndSwapInt32(&b.state, stateLocked, stateRunning)
			return nil
		}
	}
}

func (b *OnDemandBlockTaskPool) goroutine(id int) {

	// 刚启动的协程除非恰巧赶上Shutdown/ShutdownNow被调用，否则应该至少执行一个task
	idleTimer := time.NewTimer(0)
	if !idleTimer.Stop() {
		<-idleTimer.C
	}

	for {
		// log.Println("id", id, "working for loop")
		select {
		case <-b.shutdownNowCtx.Done():
			// log.Printf("id %d shutdownNow, timeoutGroup.Size=%d left\n", id, b.timeoutGroup.size())
			b.decreaseTotalGo(1)
			return
		case <-idleTimer.C:
			b.mutex.Lock()
			b.totalGo--
			b.timeoutGroup.delete(id)
			// log.Printf("id %d timeout, timeoutGroup.Size=%d left\n", id, b.timeoutGroup.size())
			b.mutex.Unlock()
			return
		case task, ok := <-b.queue:
			// log.Println("id", id, "running tasks")
			if b.timeoutGroup.isIn(id) {
				// timer只保证至少在等待X时间后才发送信号而不是在X时间内发送信号
				b.timeoutGroup.delete(id)
				// timer的Stop方法不保证一定成功
				// 不加判断并将信号清除可能会导致协程下次在case<-idleTimer.C处退出
				if !idleTimer.Stop() {
					<-idleTimer.C
				}
				// log.Println("id", id, "out timeoutGroup")
			}
			atomic.AddInt32(&b.numGoRunningTasks, 1)
			if !ok {
				// b.numGoRunningTasks > 1表示虽然当前协程监听到了b.queue关闭但还有其他协程运行task，当前协程自己退出就好
				// b.numGoRunningTasks == 1表示只有当前协程"运行task"中，其他协程在一定在"拿到b.queue到已关闭"，这一信号的路上
				// 绝不会处于运行task中
				if atomic.LoadInt32(&b.state) == stateClosing && atomic.CompareAndSwapInt32(&b.numGoRunningTasks, 1, 0) {
					// 在b.queue关闭后，第一个检测到全部task已经自然结束的协程
					b.shutdownOnce.Do(func() {
						// 状态迁移
						atomic.CompareAndSwapInt32(&b.state, stateClosing, stateStopped)
						// 显示通知外部调用者
						b.shutdownDone <- struct{}{}
						close(b.shutdownDone)
					})

					b.decreaseTotalGo(1)
					return
				}

				// 有其他协程运行task中，自己退出就好。
				atomic.AddInt32(&b.numGoRunningTasks, -1)
				b.decreaseTotalGo(1)
				return
			}
			// todo handle error
			_ = task.Run(b.shutdownNowCtx)
			atomic.AddInt32(&b.numGoRunningTasks, -1)

			b.mutex.Lock()
			// log.Println("id", id, "totalGo-mem", b.totalGo-b.timeoutGroup.size(), "totalGo", b.totalGo, "mem", b.timeoutGroup.size())
			if b.coreGo < b.totalGo && (len(b.queue) == 0 || int32(len(b.queue)) < b.totalGo) {
				// 协程在(coreGo,maxGo]区间
				// 如果没有任务可以执行，或者被判定为可能抢不到任务的协程直接退出
				// 注意：一定要在此处减1才能让此刻等待在mutex上的其他协程被正确地分区
				b.totalGo--
				// log.Println("id", id, "exits....")
				b.mutex.Unlock()
				return
			}

			if b.initGo < b.totalGo-b.timeoutGroup.size() /* && len(b.queue) == 0 */ {
				// log.Println("id", id, "initGo", b.initGo, "totalGo-mem", b.totalGo-b.timeoutGroup.size(), "totalGo", b.totalGo)
				// 协程在(initGo，coreGo]区间，如果没有任务可以执行，重置计时器
				// 当len(b.queue) != 0时，即便协程属于(coreGo,maxGo]区间，也应该给它一个定时器兜底。
				// 因为现在看队列中有任务，等真去拿的时候可能恰好没任务，如果不给它一个定时器兜底此时就会出现当前协程总数长时间大于始协程数（initGo）的情况。
				// 直到队列再次有任务时才可能将当前总协程数准确无误地降至初始协程数，因此注释掉len(b.queue) == 0判断条件
				idleTimer = time.NewTimer(b.maxIdleTime)
				b.timeoutGroup.add(id)
				// log.Println("id", id, "add timeoutGroup", "size", b.timeoutGroup.size())
			}

			b.mutex.Unlock()
		}
	}
}

func (b *OnDemandBlockTaskPool) increaseTotalGo(n int32) {
	b.mutex.Lock()
	b.totalGo += n
	b.mutex.Unlock()
}

func (b *OnDemandBlockTaskPool) decreaseTotalGo(n int32) {
	b.mutex.Lock()
	b.totalGo -= n
	b.mutex.Unlock()
}

// Shutdown 将会拒绝提交新的任务，但是会继续执行已提交任务
// 当执行完毕后，会往返回的 chan 中丢入信号
// Shutdown 会负责关闭返回的 chan
// Shutdown 无法中断正在执行的任务
func (b *OnDemandBlockTaskPool) Shutdown() (<-chan struct{}, error) {

	for {

		if atomic.LoadInt32(&b.state) == stateCreated {
			return nil, fmt.Errorf("%w", errTaskPoolIsNotRunning)
		}

		if atomic.LoadInt32(&b.state) == stateStopped {
			return nil, fmt.Errorf("%w", errTaskPoolIsStopped)
		}

		if atomic.LoadInt32(&b.state) == stateClosing {
			return nil, fmt.Errorf("%w", errTaskPoolIsClosing)
		}

		if atomic.CompareAndSwapInt32(&b.state, stateRunning, stateClosing) {
			// 目标：不但希望正在运行中的任务自然退出，还希望队列中等待的任务也能启动执行并自然退出
			// 策略：先将队列中的任务启动并执行（清空队列），再等待全部运行中的任务自然退出

			// 先关闭等待队列不再允许提交
			// 同时工作协程能够通过判断b.queue是否被关闭来终止获取任务循环
			close(b.queue)
			return b.shutdownDone, nil
		}

	}
}

// ShutdownNow 立刻关闭任务池，并且返回所有剩余未执行的任务（不包含正在执行的任务）
func (b *OnDemandBlockTaskPool) ShutdownNow() ([]Task, error) {

	for {

		if atomic.LoadInt32(&b.state) == stateCreated {
			return nil, fmt.Errorf("%w", errTaskPoolIsNotRunning)
		}

		if atomic.LoadInt32(&b.state) == stateClosing {
			return nil, fmt.Errorf("%w", errTaskPoolIsClosing)
		}

		if atomic.LoadInt32(&b.state) == stateStopped {
			return nil, fmt.Errorf("%w", errTaskPoolIsStopped)
		}

		if atomic.CompareAndSwapInt32(&b.state, stateRunning, stateStopped) {
			// 目标：立刻关闭并且返回所有剩下未执行的任务
			// 策略：关闭等待队列不再接受新任务，中断工作协程的获取任务循环，清空等待队列并保存返回

			close(b.queue)

			// 发送中断信号，中断工作协程获取任务循环
			b.shutdownNowCancel()

			// 清空队列并保存
			tasks := make([]Task, 0, len(b.queue))
			for task := range b.queue {
				tasks = append(tasks, task)
			}
			return tasks, nil
		}
	}
}

// internalState 用于查看TaskPool状态
func (b *OnDemandBlockTaskPool) internalState() int32 {
	for {
		state := atomic.LoadInt32(&b.state)
		if state != stateLocked {
			return state
		}
	}
}

// numOfGo 用于查看TaskPool中有多少工作协程
func (b *OnDemandBlockTaskPool) numOfGo() int32 {
	var n int32
	b.mutex.RLock()
	n = b.totalGo
	b.mutex.RUnlock()
	return n
}

func (b *OnDemandBlockTaskPool) States(ctx context.Context, interval time.Duration) (<-chan State, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if b.shutdownNowCtx.Err() != nil {
		return nil, b.shutdownNowCtx.Err()
	}

	statsChan := make(chan State)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case timeStamp := <-ticker.C:
				b.sendState(statsChan, timeStamp.UnixNano())
			case <-ctx.Done():
				b.sendState(statsChan, time.Now().UnixNano())
				close(statsChan)
				return
			case <-b.shutdownNowCtx.Done():
				b.sendState(statsChan, time.Now().UnixNano())
				close(statsChan)
				return
			}
		}
	}()
	return statsChan, nil
}

func (b *OnDemandBlockTaskPool) sendState(ch chan<- State, timeStamp int64) {
	// 这里发送 State 不成功则直接丢弃，不考虑重试逻辑，用户对自己的行为负责
	select {
	case ch <- b.getState(timeStamp):
	default:
	}
}

func (b *OnDemandBlockTaskPool) getState(timeStamp int64) State {
	s := State{
		PoolState:       atomic.LoadInt32(&b.state),
		GoCnt:           b.numOfGo(),
		QueueSize:       cap(b.queue),
		WaitingTasksCnt: len(b.queue),
		RunningTasksCnt: atomic.LoadInt32(&b.numGoRunningTasks),
		Timestamp:       timeStamp,
	}
	return s
}
