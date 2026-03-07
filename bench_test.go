package collections

import (
	"fmt"
	"sync"
	"testing"
)

// --- SyncMap benchmarks ---

func BenchmarkSyncMapStore(b *testing.B) {
	m := NewSyncMap[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Store(i, i)
	}
}

func BenchmarkSyncMapLoad(b *testing.B) {
	m := NewSyncMap[int, int]()
	for i := 0; i < 1000; i++ {
		m.Store(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Load(i % 1000)
	}
}

func BenchmarkSyncMapLoadOrStore(b *testing.B) {
	m := NewSyncMap[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.LoadOrStore(i%100, i)
	}
}

func BenchmarkSyncMapDelete(b *testing.B) {
	m := NewSyncMap[int, int]()
	for i := 0; i < b.N; i++ {
		m.Store(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Delete(i)
	}
}

func BenchmarkSyncMapRange(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			m := NewSyncMap[int, int]()
			for i := 0; i < size; i++ {
				m.Store(i, i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Range(func(k, v int) bool { return true })
			}
		})
	}
}

func BenchmarkSyncMapConcurrentReadWrite(b *testing.B) {
	m := NewSyncMap[int, int]()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				m.Store(i%100, i)
			} else {
				m.Load(i % 100)
			}
			i++
		}
	})
}

// --- Queue benchmarks ---

func BenchmarkQueuePush(b *testing.B) {
	q := NewQueue[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
}

func BenchmarkQueuePushPop(b *testing.B) {
	q := NewQueue[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(i)
		q.Pop()
	}
}

func BenchmarkQueueDrain(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				q := NewQueue[int]()
				for j := 0; j < size; j++ {
					q.Push(j)
				}
				b.StartTimer()
				q.Drain()
			}
		})
	}
}

func BenchmarkBoundedQueuePush(b *testing.B) {
	q := NewBoundedQueue[int](1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(i)
		if q.IsFull() {
			q.Pop()
		}
	}
}

func BenchmarkQueueConcurrentPushPop(b *testing.B) {
	q := NewQueue[int]()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				q.Push(i)
			} else {
				q.Pop()
			}
			i++
		}
	})
}

// --- Set benchmarks ---

func BenchmarkSetAdd(b *testing.B) {
	s := NewSet[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}

func BenchmarkSetContains(b *testing.B) {
	s := NewSet[int]()
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Contains(i % 1000)
	}
}

func BenchmarkSetRemove(b *testing.B) {
	s := NewSet[int]()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Remove(i)
	}
}

func BenchmarkSetUnion(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			a := NewSet[int]()
			c := NewSet[int]()
			for i := 0; i < size; i++ {
				a.Add(i)
				c.Add(i + size/2)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				a.Union(c)
			}
		})
	}
}

func BenchmarkSetIntersection(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			a := NewSet[int]()
			c := NewSet[int]()
			for i := 0; i < size; i++ {
				a.Add(i)
				c.Add(i + size/2)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				a.Intersection(c)
			}
		})
	}
}

func BenchmarkSetDifference(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			a := NewSet[int]()
			c := NewSet[int]()
			for i := 0; i < size; i++ {
				a.Add(i)
				c.Add(i + size/2)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				a.Difference(c)
			}
		})
	}
}

func BenchmarkSetConcurrentAddContains(b *testing.B) {
	s := NewSet[int]()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%4 == 0 {
				s.Add(i % 100)
			} else {
				s.Contains(i % 100)
			}
			i++
		}
	})
}

// --- RingBuffer benchmarks ---

func BenchmarkRingBufferPush(b *testing.B) {
	for _, cap := range []int{64, 1024, 8192} {
		b.Run(fmt.Sprintf("cap=%d", cap), func(b *testing.B) {
			rb := NewRingBuffer[int](cap)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				rb.Push(i)
			}
		})
	}
}

func BenchmarkRingBufferPushPop(b *testing.B) {
	rb := NewRingBuffer[int](1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Push(i)
		rb.Pop()
	}
}

func BenchmarkRingBufferToSlice(b *testing.B) {
	for _, cap := range []int{64, 1024, 8192} {
		b.Run(fmt.Sprintf("cap=%d", cap), func(b *testing.B) {
			rb := NewRingBuffer[int](cap)
			for i := 0; i < cap; i++ {
				rb.Push(i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				rb.ToSlice()
			}
		})
	}
}

func BenchmarkRingBufferOverwrite(b *testing.B) {
	rb := NewRingBuffer[int](64)
	// Fill to capacity first
	for i := 0; i < 64; i++ {
		rb.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Push(i)
	}
}

func BenchmarkRingBufferConcurrentPushPop(b *testing.B) {
	rb := NewRingBuffer[int](1024)
	// Pre-fill half so pops succeed
	for i := 0; i < 512; i++ {
		rb.Push(i)
	}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				rb.Push(i)
			} else {
				rb.Pop()
			}
			i++
		}
	})
}

// --- Pool benchmarks ---

func BenchmarkPoolGetPut(b *testing.B) {
	p := NewPool(func() []byte { return make([]byte, 1024) }, 100)
	p.Prefill(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := p.Get()
		p.Put(buf)
	}
}

func BenchmarkPoolGetFactory(b *testing.B) {
	p := NewPool(func() []byte { return make([]byte, 1024) }, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Get()
	}
}

func BenchmarkPoolConcurrentGetPut(b *testing.B) {
	p := NewPool(func() []byte { return make([]byte, 1024) }, 0)
	p.Prefill(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := p.Get()
			p.Put(buf)
		}
	})
}

func BenchmarkPoolVsSyncPool(b *testing.B) {
	factory := func() []byte { return make([]byte, 1024) }

	b.Run("Pool", func(b *testing.B) {
		p := NewPool(factory, 0)
		p.Prefill(100)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf := p.Get()
			p.Put(buf)
		}
	})

	b.Run("sync.Pool", func(b *testing.B) {
		p := &sync.Pool{New: func() any { return make([]byte, 1024) }}
		// Pre-fill
		for i := 0; i < 100; i++ {
			p.Put(make([]byte, 1024))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf := p.Get().([]byte)
			p.Put(buf)
		}
	})
}

// --- LRU benchmarks ---

func BenchmarkLRUPut(b *testing.B) {
	c := NewLRU[int, int](1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Put(i, i)
	}
}

func BenchmarkLRUGet(b *testing.B) {
	c := NewLRU[int, int](1000)
	for i := 0; i < 1000; i++ {
		c.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(i % 1000)
	}
}

func BenchmarkLRUPutEvict(b *testing.B) {
	c := NewLRU[int, int](100)
	// Fill to capacity
	for i := 0; i < 100; i++ {
		c.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Put(i+100, i)
	}
}

func BenchmarkLRUConcurrentGetPut(b *testing.B) {
	c := NewLRU[int, int](1000)
	for i := 0; i < 1000; i++ {
		c.Put(i, i)
	}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				c.Put(i%1000, i)
			} else {
				c.Get(i % 1000)
			}
			i++
		}
	})
}

// --- PriorityQueue benchmarks ---

func BenchmarkPQPush(b *testing.B) {
	pq := NewPriorityQueue(func(a, b int) bool { return a < b })
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pq.Push(i)
	}
}

func BenchmarkPQPushPop(b *testing.B) {
	pq := NewPriorityQueue(func(a, b int) bool { return a < b })
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pq.Push(i)
		pq.Pop()
	}
}

func BenchmarkPQDrain(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				pq := NewPriorityQueue(func(a, b int) bool { return a < b })
				for j := 0; j < size; j++ {
					pq.Push(j)
				}
				b.StartTimer()
				pq.Drain()
			}
		})
	}
}

func BenchmarkPQConcurrentPushPop(b *testing.B) {
	pq := NewPriorityQueue(func(a, b int) bool { return a < b })
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				pq.Push(i)
			} else {
				pq.Pop()
			}
			i++
		}
	})
}
