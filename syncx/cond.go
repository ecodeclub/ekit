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
	"sync"
	"sync/atomic"
	"unsafe"
)

// Cond 实现了一个条件变量，是等待或宣布一个事件发生的goroutines的汇合点。
//
// 在改变条件和调用Wait方法的时候，Cond 关联的锁对象 L (*Mutex 或者 *RWMutex)必须被加锁,
//
// 在Go内存模型的术语中，Cond 保证 Broadcast或Signal的调用 同步于 因此而解除阻塞的 Wait 之前。
//
// 绝大多数简单用例, 最好使用 channels 而不是 Cond
// (Broadcast 对应于关闭一个 channel, Signal 对应于给一个 channel 发送消息).
type Cond struct {
	noCopy noCopy
	// L 在观察或改变条件时被加锁
	L          sync.Locker
	notifyList *notifyList
	// 用于指向自身的指针，可以用于检测是否被复制使用
	checker unsafe.Pointer
	// 用于初始化notifyList
	once sync.Once
}

// NewCond 返回 关联了 l 的新 Cond .
func NewCond(l sync.Locker) *Cond {
	return &Cond{L: l}
}

// Wait 自动解锁 c.L 并挂起当前调用的 goroutine. 在恢复执行之后 Wait 在返回前将加 c.L 锁成功.
// 和其它系统不一样, 除非调用 Broadcast 或 Signal 或者 ctx 超时了，否则 Wait 不会返回.
//
// 成功唤醒时, 返回 nil. 超时失败时, 返回ctx.Err().
// 如果 ctx 超时了, Wait 可能依旧执行成功返回 nil.
//
// 在 Wait 第一次继续执行时，因为 c.L 没有加锁, 当 Wait 返回的时候，调用者通常不能假设条件是真的
// 相反, caller 应该在循环中调用 Wait:
//
//		c.L.Lock()
//		for !condition() {
//		    if err := c.Wait(ctx); err != nil {
//	          // 超时唤醒了，并不是被正常唤醒的，可以做一些超时的处理
//			}
//		}
//		... condition 满足了，do work ...
//		c.L.Unlock()
func (c *Cond) Wait(ctx context.Context) error {
	c.checkCopy()
	c.checkFirstUse()
	t := c.notifyList.add() // 解锁前，将等待的对象放入链表中
	c.L.Unlock()            // 一定是在等待对象放入链表后再解锁，避免刚解锁就发生协程切换，执行了signal后，再换回来导致永远阻塞
	defer c.L.Lock()
	return c.notifyList.wait(ctx, t)
}

// Signal 唤醒一个等待在 c 上的goroutine.
//
// 调用时，caller 可以持有也可以不持有 c.L 锁
//
// Signal() 不影响 goroutine 调度的优先级; 如果其它的 goroutines
// 尝试着锁定 c.L, 它们可能在 "waiting" goroutine 之前被唤醒.
func (c *Cond) Signal() {
	c.checkCopy()
	c.checkFirstUse()
	c.notifyList.notifyOne()
}

// Broadcast 唤醒所有等待在 c 上的goroutine.
//
// 调用时，caller 可以持有也可以不持有 c.L 锁
func (c *Cond) Broadcast() {
	c.checkCopy()
	c.checkFirstUse()
	c.notifyList.notifyAll()
}

// checkCopy 检查是否被拷贝使用
func (c *Cond) checkCopy() {
	// 判断checker保存的指针是否等于当前的指针（初始化时，并没有初始化checker的值，所以也会出现不相等）
	if c.checker != unsafe.Pointer(c) &&
		// 由于初次初始化时，c.checker为0值，所以顺便进行一次原子替换，辅助初始化
		!atomic.CompareAndSwapPointer(&c.checker, nil, unsafe.Pointer(c)) &&
		// 再度检查checker保留指针是否等于当前指针
		c.checker != unsafe.Pointer(c) {
		panic("syncx.Cond is copied")
	}
}

// checkFirstUse 用于初始化notifyList
func (c *Cond) checkFirstUse() {
	c.once.Do(func() {
		if c.notifyList == nil {
			c.notifyList = newNotifyList()
		}
	})
}

// notifyList 是一个简单的 runtime_notifyList 实现，但增强了 wait 方法
type notifyList struct {
	mu   sync.Mutex
	list *chanList
}

func newNotifyList() *notifyList {
	return &notifyList{
		mu:   sync.Mutex{},
		list: newChanList(),
	}
}

func (l *notifyList) add() *node {
	l.mu.Lock()
	defer l.mu.Unlock()
	el := l.list.alloc()
	l.list.pushBack(el)
	return el
}

func (l *notifyList) wait(ctx context.Context, elem *node) error {
	ch := elem.Value
	// 回收ch，超时时，因为没有被使用过，直接复用
	// 正常唤醒时，由于被放入了一条消息，但被取出来了一次，所以elem中的ch可以重复使用
	// 由于ch是挂在elem上的，所以elem在ch被回收之前，不可以被错误回收，所以必须在这里进行回收
	defer l.list.free(elem)
	select { // 由于会随机选择一条，在超时和通知同时存在的话，如果通知先行，则没有影响，如果超时的同时，又来了通知
	case <-ctx.Done(): // 进了超时分支
		l.mu.Lock()
		defer l.mu.Unlock()
		select {
		// double check: 检查是否在加锁前，刚好被正常通知了，
		// 如果取到数据，代表收到了信号了，ch也因为被取了一次消息，意味着可以再次复用
		// 转移信号到下一个
		// 如果有下一个等待的，就唤醒它
		case <-ch:
			if l.list.len() != 0 {
				l.notifyNext()
			}
		// 如果取不到数据，代表不可能被正常唤醒了，ch也意味着没有被使用，可以从队列移除等待对象
		default:
			l.list.remove(elem)
		}
		return ctx.Err()
	case <-ch: // 如果取到数据，代表被正常唤醒了，ch也因为被取了一次消息，意味着可以再次复用
		return nil
	}
}

func (l *notifyList) notifyOne() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list.len() == 0 {
		return
	}
	l.notifyNext()
}

func (l *notifyList) notifyNext() {
	front := l.list.front()
	ch := front.Value
	l.list.remove(front)
	ch <- struct{}{}
}

func (l *notifyList) notifyAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for l.list.len() != 0 {
		l.notifyNext()
	}
}

// node 保存chan的链表元素
type node struct {
	prev  *node
	next  *node
	Value chan struct{}
}

// chanList 用于存放保存channel的一个双链表， 带复用元素的功能
type chanList struct {
	// 哨兵元素，方便处理元素个数为0的情况
	sentinel *node
	size     int
	pool     *sync.Pool
}

func newChanList() *chanList {
	sentinel := &node{}
	sentinel.prev = sentinel
	sentinel.next = sentinel
	return &chanList{
		sentinel: sentinel,
		size:     0,
		pool: &sync.Pool{
			New: func() any {
				return &node{
					Value: make(chan struct{}, 1),
				}
			},
		},
	}
}

// len 获取链表长度
func (l *chanList) len() int {
	return l.size
}

// front 获取队首元素
func (l *chanList) front() *node {
	return l.sentinel.next
}

// alloc 申请新的元素，包含复用的chan
func (l *chanList) alloc() *node {
	elem := l.pool.Get().(*node)
	return elem
}

// pushBack 追加元素到队尾
func (l *chanList) pushBack(elem *node) {
	elem.next = l.sentinel
	elem.prev = l.sentinel.prev
	l.sentinel.prev.next = elem
	l.sentinel.prev = elem
	l.size++
}

// remove 元素移除时，还不能回收该元素，避免元素上的chan被错误覆盖
func (l *chanList) remove(elem *node) {
	elem.prev.next = elem.next
	elem.next.prev = elem.prev
	elem.prev = nil
	elem.next = nil
	l.size--
}

// free 回收该元素，用于下次alloc获取时复用，避免再次分配
func (l *chanList) free(elem *node) {
	l.pool.Put(elem)
}

// 用于静态代码检查复制的问题
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
