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
