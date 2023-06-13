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

type notifyList struct {
	mu   sync.Mutex
	list *list.List
}

func newNotifyList() *notifyList {
	return &notifyList{
		mu:   sync.Mutex{},
		list: list.New(),
	}
}

func (l *notifyList) add() *list.Element {
	l.mu.Lock()
	defer l.mu.Unlock()
	ch := make(chan struct{})
	return l.list.PushBack(ch)
}

func (l *notifyList) wait(elem *list.Element) {
	ch := elem.Value.(chan struct{})
	<-ch
}

func (l *notifyList) waitWithContext(ctx context.Context, elem *list.Element) error {
	ch := elem.Value.(chan struct{})
	select { // 由于会随机选择一条，在超时和通知同时存在的话，如果通知先行，则没有影响，如果超时的同时，又来了通知
	case <-ctx.Done(): // 进了超时分支，但同时协程发生了切换进入了notifyOne的分支；这个时候，根据remove的成功与否可以知道是否是需要唤醒的
		l.mu.Lock()
		defer l.mu.Unlock()
		select {
		// double check: 检查是否在加锁前，刚好被正常通知了，
		// 这种情况应该是正常消费的情况，等同于在恰巧超时时刻被唤醒，修正成正常唤醒的情况
		case <-ch:
			return nil
		default:
			// 这种情况代表加锁成功后，没有被通知到，属于真正的超时的情况，从队列移除等待对象，避免被错误通知唤醒，返回超时错误信息
			l.list.Remove(elem)
			close(ch)
			return ctx.Err()
		}
	case <-ch:
		return nil
	}
}

func (l *notifyList) notifyOne() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list.Len() == 0 {
		return
	}
	ch := l.list.Front().Value.(chan struct{})
	l.list.Remove(l.list.Front())
	close(ch)
}

func (l *notifyList) notifyAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for l.list.Len() != 0 {
		ch := l.list.Front().Value.(chan struct{})
		l.list.Remove(l.list.Front())
		close(ch)
	}
}

type Cond struct {
	L          sync.Locker
	notifyList *notifyList
}

func NewCond(l sync.Locker) *Cond {
	return &Cond{
		L:          l,
		notifyList: newNotifyList(),
	}
}

func (c *Cond) Wait() {
	t := c.notifyList.add() // 解锁前，将等待的对象放入链表中
	c.L.Unlock()            // 一定是在等待对象放入链表后再解锁，避免刚解锁就发生协程切换，执行了signal后，再换回来导致永远阻塞
	defer c.L.Lock()
	c.notifyList.wait(t)
}

func (c *Cond) WaitWithContext(ctx context.Context) error {
	t := c.notifyList.add()
	c.L.Unlock()
	defer c.L.Lock()
	return c.notifyList.waitWithContext(ctx, t)
}

func (c *Cond) Signal() {
	c.notifyList.notifyOne()
}

func (c *Cond) Broadcast() {
	c.notifyList.notifyAll()
}
