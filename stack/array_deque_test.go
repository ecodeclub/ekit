package stack

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

type person struct {
	name string
	age  int
}

func TestArrayDeque_Push1(t *testing.T) {
	deque := NewArrayDeque[int]()
	err := deque.Push(1)
	if err != nil {
		t.Logf("err: %v", err)
	}
	pop, err := deque.Pop()
	if err != nil {
		t.Logf("err: %v", err)
	}
	t.Logf("pop value: %v", pop)
}

func TestArrayDeque_Push2(t *testing.T) {
	deque := NewArrayDeque[int]()
	for i := 0; i < 63; i++ {
		err := deque.Push(i)
		if err != nil {
			t.Logf("push err: %v", err)
		}
	}
	for !deque.IsEmpty() {
		pop, err := deque.Pop()
		if err != nil {
			t.Logf("err: %v", err)
		}
		t.Logf("pop value: %v", pop)
	}
}

func TestArrayDeque_Push3(t *testing.T) {
	clazzList := []*person{
		&person{"zhangsan", 1},
		&person{"lisi", 2},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
		&person{"wangwu", 3},
	}
	deque := NewArrayDeque[*person]()
	for _, val := range clazzList {
		err := deque.Push(val)
		if err != nil {
			t.Logf("err: %v", err)
		}
	}
	for !deque.IsEmpty() {
		pop, err := deque.Pop()
		if err != nil {
			t.Logf("err: %v", err)
		}
		t.Logf("pop value: %v", pop)
	}
}

func TestArrayDeque_Push(t *testing.T) {
	stack := list.New()
	deque := NewArrayDeque[int]()
	n := 100000
	for i := 0; i < n; i++ {
		stack.PushBack(i)
		err := deque.Push(i)
		assert.Nil(t, err, "must return nil")
	}
	for stack.Len() > 0 {
		ele1 := stack.Remove(stack.Back()).(int)
		pop, err := deque.Pop()
		assert.Nil(t, err, "error must be nil")
		assert.Equal(t, ele1, pop, "")
	}
}
