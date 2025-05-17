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

import "testing"

func BenchmarkPriorityQueue_NoCapacity_Enqueue(b *testing.B) {
	n := 300000
	pq := priorityQueueOf(-1, []int{}, compare())
	b.ResetTimer()
	for i := n; i > 0; i-- {
		if pq.Enqueue(i) != nil {
			b.Fail()
		}
	}
}

func BenchmarkPriorityQueue_NoCapacity_Dequeue(b *testing.B) {
	n := 300000
	pq := priorityQueueOf(-1, []int{}, compare())
	for i := n; i > 0; i-- {
		if pq.Enqueue(i) != nil {
			b.Fail()
		}
	}
	b.ResetTimer()
	for i := 0; i < n; i++ {
		_, err := pq.Dequeue()
		if err != nil {
			b.Fail()
		}
	}
}

func BenchmarkPriorityQueue_Capacity_Enqueue(b *testing.B) {
	n := 300000
	pq := priorityQueueOf(n, []int{}, compare())
	b.ResetTimer()
	for i := n; i > 0; i-- {
		if pq.Enqueue(i) != nil {
			b.Fail()
		}
	}
}

func BenchmarkPriorityQueue_Capacity_Dequeue(b *testing.B) {
	n := 300000
	pq := priorityQueueOf(n, []int{}, compare())
	for i := n; i > 0; i-- {
		if pq.Enqueue(i) != nil {
			b.Fail()
		}
	}
	b.ResetTimer()
	for i := 0; i < n; i++ {
		_, err := pq.Dequeue()
		if err != nil {
			b.Fail()
		}
	}
}
