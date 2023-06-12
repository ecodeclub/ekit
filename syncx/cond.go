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

type notifyItem struct {
	list *notifyList
	ch   chan struct{}
	elem *list.Element
}

func newNotifyItem(l *notifyList) *notifyItem {
	return &notifyItem{list: l, ch: make(chan struct{})}
}

func (n *notifyItem) notify() {
	close(n.ch)
}

func (n *notifyItem) wait() {
	<-n.ch
}

func (n *notifyItem) waitWithContext(ctx context.Context) error {
	select { // 由于会随机选择一条，在超时和通知同时存在的话，如果通知先行，则没有影响，如果超时的同时，又来了通知
	case <-ctx.Done(): // 进了超时分支，但同时协程发生了切换进入了notifyOne的分支；这个时候，根据remove的成功与否可以知道是否是需要唤醒的
		if n.list.remove(n) {
			return ctx.Err()
		}
		return nil
	case <-n.ch:
		return nil
	}
}

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

func (l *notifyList) add() *notifyItem {
	l.mu.Lock()
	defer l.mu.Unlock()
	item := newNotifyItem(l)
	item.elem = l.list.PushBack(item)
	return item
}

func (l *notifyList) remove(item *notifyItem) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	select {
	case <-item.ch: // 检查是否在加锁前，刚好被通知了，这种情况应该是正常消费的情况，只是因为恰好超时了而已
		return false
	default:
		l.list.Remove(item.elem)
		item.notify()
		return true
	}
}

func (l *notifyList) notifyOne() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list.Len() == 0 {
		return
	}
	item := l.list.Front().Value.(*notifyItem)
	l.list.Remove(l.list.Front())
	item.notify()
}

func (l *notifyList) notifyAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for l.list.Len() != 0 {
		item := l.list.Front().Value.(*notifyItem)
		l.list.Remove(l.list.Front())
		item.notify()
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
	notifyItem := c.notifyList.add() // 解锁前，将等待的对象放入链表中
	c.L.Unlock()                     // 一定是在等待对象放入链表后再解锁，避免刚解锁就发生协程切换，执行了signal后，再换回来导致永远阻塞
	defer c.L.Lock()
	notifyItem.wait()
}

func (c *Cond) WaitWithContext(ctx context.Context) error {
	notifyItem := c.notifyList.add()
	c.L.Unlock()
	defer c.L.Lock()
	return notifyItem.waitWithContext(ctx)
}

func (c *Cond) Signal() {
	c.notifyList.notifyOne()
}

func (c *Cond) Broadcast() {
	c.notifyList.notifyAll()
}
