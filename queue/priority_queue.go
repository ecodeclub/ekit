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

package queue

import (
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/internal/queue"
)

type PriorityQueue[T any] struct {
	priorityQueue *queue.PriorityQueue[T]
}

func NewPriorityQueue[T any](capacity int, compare ekit.Comparator[T]) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{}
	pq.priorityQueue = queue.NewPriorityQueue[T](capacity, compare)
	return pq
}

func (pq *PriorityQueue[T]) Len() int {
	return pq.priorityQueue.Len()
}

func (pq *PriorityQueue[T]) Peek() (T, error) {
	return pq.priorityQueue.Peek()
}

func (pq *PriorityQueue[T]) Enqueue(t T) error {
	return pq.priorityQueue.Enqueue(t)
}

func (pq *PriorityQueue[T]) Dequeue() (T, error) {
	return pq.priorityQueue.Dequeue()
}
