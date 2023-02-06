package stack

import "github.com/gotomicro/ekit/queue"

type Stack[T any] interface {
	// Push 将元素入栈，如果无法入栈，那么返回错误
	Push(t T) error
	// Pop 弹出栈顶元素, 如果此时栈中没有数据，那么返回错误
	Pop() (T, error)
	// Top 查看栈顶元素，但是不删除，如果栈中没有数据，那么返回错误
	Top() (T, error)
}

type Queue[T any] interface {
	// Front 查看队首元素，但是不删除，如果队列中没有数据，那么返回错误
	Front() (T, error)
	queue.Queue[T]
}
