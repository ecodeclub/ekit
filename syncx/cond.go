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
	"container/list"
	"context"
	"sync"
)

// notifyList 是一个简单的 runtime_notifyList 实现，但增加了 waitWithContext 方法
type notifyList struct {
	mu     sync.Mutex
	list   *list.List
	chPool *sync.Pool
}

func newNotifyList() *notifyList {
	return &notifyList{
		mu:   sync.Mutex{},
		list: list.New(),
		chPool: &sync.Pool{
			New: func() any {
				return make(chan struct{}, 1)
			},
		},
	}
}

func (l *notifyList) add() *list.Element {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.list.PushBack(l.chPool.Get())
}

func (l *notifyList) waitWithContext(ctx context.Context, elem *list.Element) error {
	ch := elem.Value.(chan struct{})
	// 回收ch，超时时，因为没有被使用过，直接复用
	// 正常唤醒时，由于被放入了一条消息，但被取出来了一次，所以可以重复使用
	defer l.chPool.Put(ch)
	select { // 由于会随机选择一条，在超时和通知同时存在的话，如果通知先行，则没有影响，如果超时的同时，又来了通知
	case <-ctx.Done(): // 进了超时分支，但同时协程发生了切换进入了notifyOne的分支；这个时候，根据remove的成功与否可以知道是否是需要唤醒的
		l.mu.Lock()
		defer l.mu.Unlock()
		select {
		// double check: 检查是否在加锁前，刚好被正常通知了，
		// 这种情况应该是正常消费的情况，等同于在恰巧超时时刻被唤醒，修正成正常唤醒的情况
		case <-ch: // 如果取到数据，代表收到了信号了，ch也因为被取了一次消息，意味着可以再次复用
			// 转移信号到下一个
			// 如果没有下一个等待的，就返回
			if l.list.Len() == 0 {
				return ctx.Err()
			}
			// 如果有下一个等待的，就唤醒它
			l.notifyNext()
		default: // 如果取不到数据，代表不可能被正常唤醒了，ch也意味着没有被使用
			// 这种情况代表加锁成功后，没有被通知到，属于真正的超时的情况，从队列移除等待对象，避免被错误通知唤醒，返回超时错误信息
			l.list.Remove(elem)
		}
		return ctx.Err()
	case <-ch: // 如果取到数据，代表被正常唤醒了，ch也因为被取了一次消息，意味着可以再次复用
		return nil
	}
}

func (l *notifyList) notifyOne() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list.Len() == 0 {
		return
	}
	l.notifyNext()
}

func (l *notifyList) notifyNext() {
	front := l.list.Front()
	ch := front.Value.(chan struct{})
	l.list.Remove(front)
	ch <- struct{}{}
}

func (l *notifyList) notifyAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for l.list.Len() != 0 {
		l.notifyNext()
	}
}

// Cond 实现了一个条件变量，是等待或宣布一个事件发生的goroutines的汇合点。
//
// 在改变条件和调用Wait方法的时候，Cond 关联的锁对象 L (*Mutex 或者 *RWMutex)必须被加锁,
//
// # Cond 在初次使用后，不要复制对象
//
// 在Go内存模型的术语中，Cond安排对Broadcast或Signal的调用"happens before"任何解除阻塞的 Wait 调用。
//
// 绝大多数简单用例, 最好使用channels而不是 Cond
// (Broadcast 对应于关闭一个 channel, Signal 对应于给一个 channel 发送消息).
type Cond struct {
	// L 在观察或改变条件时被加锁
	L          sync.Locker
	notifyList *notifyList
}

// NewCond 返回 关联了 l 的新 Cond .
func NewCond(l sync.Locker) *Cond {
	return &Cond{
		L:          l,
		notifyList: newNotifyList(),
	}
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
	t := c.notifyList.add() // 解锁前，将等待的对象放入链表中
	c.L.Unlock()            // 一定是在等待对象放入链表后再解锁，避免刚解锁就发生协程切换，执行了signal后，再换回来导致永远阻塞
	defer c.L.Lock()
	return c.notifyList.waitWithContext(ctx, t)
}

// Signal 唤醒一个等待在 c 上的goroutine.
//
// 调用时，caller 可以持有也可以不持有 c.L 锁
//
// Signal() 不影响 goroutine 调度的优先级; 如果其它的 goroutines
// 尝试着锁定 c.L, 它们可能在 "waiting" goroutine 之前被唤醒.
func (c *Cond) Signal() {
	c.notifyList.notifyOne()
}

// Broadcast 唤醒所有等待在 c 上的goroutine.
//
// 调用时，caller 可以持有也可以不持有 c.L 锁
func (c *Cond) Broadcast() {
	c.notifyList.notifyAll()
}
