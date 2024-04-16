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

package queue_test

import (
	"fmt"

	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/internal/queue"
)

func ExampleNewPriorityQueue() {
	// 容量，并且队列里面放的是 int
	pq := queue.NewPriorityQueue(10, ekit.ComparatorRealNumber[int])
	_ = pq.Enqueue(10)
	_ = pq.Enqueue(9)
	val, _ := pq.Dequeue()
	fmt.Println(val)
	// Output:
	// 9
}
