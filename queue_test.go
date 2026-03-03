package collections

import (
	"sync"
	"testing"
)

func TestQueue_PushPop(t *testing.T) {
	q := NewQueue[int]()
	q.Push(1)
	q.Push(2)
	q.Push(3)

	for _, want := range []int{1, 2, 3} {
		got, ok := q.Pop()
		if !ok || got != want {
			t.Fatalf("Pop = %d, %v; want %d, true", got, ok, want)
		}
	}
}

func TestQueue_PopEmpty(t *testing.T) {
	q := NewQueue[string]()
	v, ok := q.Pop()
	if ok || v != "" {
		t.Fatalf("Pop(empty) = %q, %v; want '', false", v, ok)
	}
}

func TestQueue_Peek(t *testing.T) {
	q := NewQueue[int]()
	_, ok := q.Peek()
	if ok {
		t.Fatal("Peek(empty) should return false")
	}

	q.Push(42)
	v, ok := q.Peek()
	if !ok || v != 42 {
		t.Fatalf("Peek = %d, %v; want 42, true", v, ok)
	}

	// Peek should not remove.
	if q.Len() != 1 {
		t.Fatal("Peek should not remove items")
	}
}

func TestQueue_Len(t *testing.T) {
	q := NewQueue[int]()
	if q.Len() != 0 {
		t.Fatalf("empty Len = %d; want 0", q.Len())
	}

	q.Push(1)
	q.Push(2)
	if q.Len() != 2 {
		t.Fatalf("Len = %d; want 2", q.Len())
	}
}

func TestQueue_IsEmpty(t *testing.T) {
	q := NewQueue[int]()
	if !q.IsEmpty() {
		t.Fatal("new queue should be empty")
	}

	q.Push(1)
	if q.IsEmpty() {
		t.Fatal("queue with items should not be empty")
	}
}

func TestQueue_Bounded(t *testing.T) {
	q := NewBoundedQueue[int](3)

	if !q.Push(1) || !q.Push(2) || !q.Push(3) {
		t.Fatal("pushes within capacity should succeed")
	}

	if q.Push(4) {
		t.Fatal("push beyond capacity should fail")
	}

	if q.Len() != 3 {
		t.Fatalf("Len = %d; want 3", q.Len())
	}
}

func TestQueue_BoundedIsFull(t *testing.T) {
	q := NewBoundedQueue[int](2)
	if q.IsFull() {
		t.Fatal("empty bounded queue should not be full")
	}

	q.Push(1)
	q.Push(2)
	if !q.IsFull() {
		t.Fatal("full bounded queue should report IsFull")
	}

	q.Pop()
	if q.IsFull() {
		t.Fatal("after pop, queue should not be full")
	}
}

func TestQueue_UnboundedIsNeverFull(t *testing.T) {
	q := NewQueue[int]()
	q.Push(1)
	if q.IsFull() {
		t.Fatal("unbounded queue should never be full")
	}
}

func TestQueue_BoundedZeroCapacity(t *testing.T) {
	q := NewBoundedQueue[int](0)
	// Should default to capacity 1.
	if !q.Push(1) {
		t.Fatal("first push should succeed")
	}
	if q.Push(2) {
		t.Fatal("second push should fail (capacity=1)")
	}
}

func TestQueue_Clear(t *testing.T) {
	q := NewQueue[int]()
	q.Push(1)
	q.Push(2)
	q.Clear()

	if !q.IsEmpty() {
		t.Fatal("queue should be empty after Clear")
	}
}

func TestQueue_Drain(t *testing.T) {
	q := NewQueue[int]()
	q.Push(10)
	q.Push(20)
	q.Push(30)

	items := q.Drain()
	if len(items) != 3 || items[0] != 10 || items[1] != 20 || items[2] != 30 {
		t.Fatalf("Drain = %v; want [10 20 30]", items)
	}

	if !q.IsEmpty() {
		t.Fatal("queue should be empty after Drain")
	}
}

func TestQueue_DrainEmpty(t *testing.T) {
	q := NewQueue[int]()
	items := q.Drain()
	if items != nil {
		t.Fatalf("Drain(empty) = %v; want nil", items)
	}
}

func TestQueue_FIFO(t *testing.T) {
	q := NewQueue[string]()
	words := []string{"alpha", "beta", "gamma", "delta"}
	for _, w := range words {
		q.Push(w)
	}
	for _, want := range words {
		got, _ := q.Pop()
		if got != want {
			t.Fatalf("Pop = %q; want %q", got, want)
		}
	}
}

func TestQueue_Concurrent(t *testing.T) {
	q := NewQueue[int]()
	var wg sync.WaitGroup

	// Push from multiple goroutines.
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			q.Push(n)
		}(i)
	}

	wg.Wait()
	if q.Len() != 100 {
		t.Fatalf("concurrent Len = %d; want 100", q.Len())
	}

	// Pop all from multiple goroutines.
	var popped sync.Map
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if v, ok := q.Pop(); ok {
				popped.Store(v, true)
			}
		}()
	}

	wg.Wait()
	if !q.IsEmpty() {
		t.Fatal("queue should be empty after concurrent drain")
	}
}
