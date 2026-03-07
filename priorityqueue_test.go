package collections

import (
	"sync"
	"testing"
)

func intLess(a, b int) bool { return a < b }

func TestPQ_PushPop(t *testing.T) {
	pq := NewPriorityQueue(intLess)
	pq.Push(3)
	pq.Push(1)
	pq.Push(2)

	for _, want := range []int{1, 2, 3} {
		got, ok := pq.Pop()
		if !ok || got != want {
			t.Fatalf("Pop = %d, %v; want %d, true", got, ok, want)
		}
	}
}

func TestPQ_PopEmpty(t *testing.T) {
	pq := NewPriorityQueue(intLess)
	v, ok := pq.Pop()
	if ok || v != 0 {
		t.Fatalf("Pop(empty) = %d, %v; want 0, false", v, ok)
	}
}

func TestPQ_Peek(t *testing.T) {
	pq := NewPriorityQueue(intLess)

	_, ok := pq.Peek()
	if ok {
		t.Fatal("Peek(empty) should return false")
	}

	pq.Push(5)
	pq.Push(2)
	pq.Push(8)

	v, ok := pq.Peek()
	if !ok || v != 2 {
		t.Fatalf("Peek = %d, %v; want 2, true", v, ok)
	}

	// Peek should not remove.
	if pq.Len() != 3 {
		t.Fatalf("Len after Peek = %d; want 3", pq.Len())
	}
}

func TestPQ_Len(t *testing.T) {
	pq := NewPriorityQueue(intLess)
	if pq.Len() != 0 {
		t.Fatalf("Len = %d; want 0", pq.Len())
	}

	pq.Push(1)
	pq.Push(2)
	if pq.Len() != 2 {
		t.Fatalf("Len = %d; want 2", pq.Len())
	}
}

func TestPQ_IsEmpty(t *testing.T) {
	pq := NewPriorityQueue(intLess)
	if !pq.IsEmpty() {
		t.Fatal("new PQ should be empty")
	}

	pq.Push(1)
	if pq.IsEmpty() {
		t.Fatal("PQ with items should not be empty")
	}
}

func TestPQ_Clear(t *testing.T) {
	pq := NewPriorityQueue(intLess)
	pq.Push(1)
	pq.Push(2)
	pq.Clear()

	if !pq.IsEmpty() {
		t.Fatal("PQ should be empty after Clear")
	}
}

func TestPQ_Drain(t *testing.T) {
	pq := NewPriorityQueue(intLess)
	pq.Push(5)
	pq.Push(1)
	pq.Push(3)

	items := pq.Drain()
	if len(items) != 3 || items[0] != 1 || items[1] != 3 || items[2] != 5 {
		t.Fatalf("Drain = %v; want [1 3 5]", items)
	}

	if !pq.IsEmpty() {
		t.Fatal("PQ should be empty after Drain")
	}
}

func TestPQ_DrainEmpty(t *testing.T) {
	pq := NewPriorityQueue(intLess)
	items := pq.Drain()
	if len(items) != 0 {
		t.Fatalf("Drain(empty) = %v; want []", items)
	}
}

func TestPQ_PushAll(t *testing.T) {
	pq := NewPriorityQueue(intLess)
	pq.PushAll(9, 3, 7, 1, 5)

	if pq.Len() != 5 {
		t.Fatalf("Len = %d; want 5", pq.Len())
	}

	v, _ := pq.Pop()
	if v != 1 {
		t.Fatalf("Pop = %d; want 1 (min)", v)
	}
}

func TestPQ_Duplicates(t *testing.T) {
	pq := NewPriorityQueue(intLess)
	pq.Push(2)
	pq.Push(2)
	pq.Push(1)
	pq.Push(2)

	items := pq.Drain()
	want := []int{1, 2, 2, 2}
	if len(items) != len(want) {
		t.Fatalf("Drain = %v; want %v", items, want)
	}
	for i, v := range items {
		if v != want[i] {
			t.Fatalf("Drain[%d] = %d; want %d", i, v, want[i])
		}
	}
}

func TestPQ_MaxHeap(t *testing.T) {
	// Use reverse comparator for max-heap behavior.
	pq := NewPriorityQueue(func(a, b int) bool { return a > b })
	pq.Push(1)
	pq.Push(5)
	pq.Push(3)

	v, _ := pq.Pop()
	if v != 5 {
		t.Fatalf("Pop = %d; want 5 (max)", v)
	}
}

type task struct {
	name     string
	priority int
}

func TestPQ_Structs(t *testing.T) {
	pq := NewPriorityQueue(func(a, b task) bool { return a.priority < b.priority })
	pq.Push(task{"low", 10})
	pq.Push(task{"high", 1})
	pq.Push(task{"med", 5})

	got, ok := pq.Pop()
	if !ok || got.name != "high" || got.priority != 1 {
		t.Fatalf("Pop = %+v; want {high 1}", got)
	}
}

func TestPQ_LargeN(t *testing.T) {
	pq := NewPriorityQueue(intLess)
	n := 10000
	for i := n; i > 0; i-- {
		pq.Push(i)
	}

	prev := 0
	for i := 0; i < n; i++ {
		v, ok := pq.Pop()
		if !ok {
			t.Fatalf("Pop failed at iteration %d", i)
		}
		if v < prev {
			t.Fatalf("out of order: %d < %d", v, prev)
		}
		prev = v
	}
}

func TestPQ_Concurrent(t *testing.T) {
	pq := NewPriorityQueue(intLess)
	var wg sync.WaitGroup

	// Concurrent pushes.
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			pq.Push(n)
		}(i)
	}

	wg.Wait()

	if pq.Len() != 200 {
		t.Fatalf("Len = %d; want 200", pq.Len())
	}

	// Concurrent pops.
	var popped sync.Map
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if v, ok := pq.Pop(); ok {
				popped.Store(v, true)
			}
		}()
	}

	wg.Wait()

	if !pq.IsEmpty() {
		t.Fatal("PQ should be empty after concurrent drain")
	}
}

func TestPQ_MixedPushPop(t *testing.T) {
	pq := NewPriorityQueue(intLess)

	pq.Push(10)
	pq.Push(5)
	v, _ := pq.Pop() // 5
	if v != 5 {
		t.Fatalf("Pop = %d; want 5", v)
	}

	pq.Push(3)
	pq.Push(7)
	v, _ = pq.Pop() // 3
	if v != 3 {
		t.Fatalf("Pop = %d; want 3", v)
	}

	v, _ = pq.Pop() // 7
	if v != 7 {
		t.Fatalf("Pop = %d; want 7", v)
	}

	v, _ = pq.Pop() // 10
	if v != 10 {
		t.Fatalf("Pop = %d; want 10", v)
	}
}

func TestPQ_StringPriority(t *testing.T) {
	pq := NewPriorityQueue(func(a, b string) bool { return a < b })
	pq.Push("cherry")
	pq.Push("apple")
	pq.Push("banana")

	items := pq.Drain()
	want := []string{"apple", "banana", "cherry"}
	for i, v := range items {
		if v != want[i] {
			t.Fatalf("Drain[%d] = %q; want %q", i, v, want[i])
		}
	}
}
