package collections

import (
	"sync"
	"testing"
)

func TestRingBuffer_PushPop(t *testing.T) {
	rb := NewRingBuffer[int](5)
	rb.Push(1)
	rb.Push(2)
	rb.Push(3)

	for _, want := range []int{1, 2, 3} {
		got, ok := rb.Pop()
		if !ok || got != want {
			t.Fatalf("Pop = %d, %v; want %d, true", got, ok, want)
		}
	}
}

func TestRingBuffer_PopEmpty(t *testing.T) {
	rb := NewRingBuffer[int](3)
	v, ok := rb.Pop()
	if ok || v != 0 {
		t.Fatalf("Pop(empty) = %d, %v; want 0, false", v, ok)
	}
}

func TestRingBuffer_Overflow(t *testing.T) {
	rb := NewRingBuffer[string](3)
	rb.Push("a")
	rb.Push("b")
	rb.Push("c")

	evicted, wasEvicted := rb.Push("d")
	if !wasEvicted || evicted != "a" {
		t.Fatalf("overflow Push = %q, %v; want 'a', true", evicted, wasEvicted)
	}

	// Buffer should now contain b, c, d.
	got, _ := rb.Pop()
	if got != "b" {
		t.Fatalf("after overflow, Pop = %q; want 'b'", got)
	}
}

func TestRingBuffer_NoEviction(t *testing.T) {
	rb := NewRingBuffer[int](5)
	evicted, wasEvicted := rb.Push(42)
	if wasEvicted || evicted != 0 {
		t.Fatalf("Push without eviction = %d, %v; want 0, false", evicted, wasEvicted)
	}
}

func TestRingBuffer_Peek(t *testing.T) {
	rb := NewRingBuffer[int](3)
	_, ok := rb.Peek()
	if ok {
		t.Fatal("Peek(empty) should return false")
	}

	rb.Push(10)
	rb.Push(20)

	v, ok := rb.Peek()
	if !ok || v != 10 {
		t.Fatalf("Peek = %d, %v; want 10, true", v, ok)
	}

	// Peek should not modify.
	if rb.Len() != 2 {
		t.Fatal("Peek should not modify buffer")
	}
}

func TestRingBuffer_PeekNewest(t *testing.T) {
	rb := NewRingBuffer[int](5)
	_, ok := rb.PeekNewest()
	if ok {
		t.Fatal("PeekNewest(empty) should return false")
	}

	rb.Push(1)
	rb.Push(2)
	rb.Push(3)

	v, ok := rb.PeekNewest()
	if !ok || v != 3 {
		t.Fatalf("PeekNewest = %d, %v; want 3, true", v, ok)
	}
}

func TestRingBuffer_PeekNewestAfterOverflow(t *testing.T) {
	rb := NewRingBuffer[int](2)
	rb.Push(1)
	rb.Push(2)
	rb.Push(3) // evicts 1

	v, ok := rb.PeekNewest()
	if !ok || v != 3 {
		t.Fatalf("PeekNewest = %d, %v; want 3, true", v, ok)
	}
}

func TestRingBuffer_Len(t *testing.T) {
	rb := NewRingBuffer[int](5)
	if rb.Len() != 0 {
		t.Fatalf("empty Len = %d; want 0", rb.Len())
	}

	rb.Push(1)
	rb.Push(2)
	if rb.Len() != 2 {
		t.Fatalf("Len = %d; want 2", rb.Len())
	}
}

func TestRingBuffer_Cap(t *testing.T) {
	rb := NewRingBuffer[int](7)
	if rb.Cap() != 7 {
		t.Fatalf("Cap = %d; want 7", rb.Cap())
	}
}

func TestRingBuffer_IsFull(t *testing.T) {
	rb := NewRingBuffer[int](2)
	if rb.IsFull() {
		t.Fatal("empty buffer should not be full")
	}

	rb.Push(1)
	rb.Push(2)
	if !rb.IsFull() {
		t.Fatal("buffer at capacity should be full")
	}
}

func TestRingBuffer_IsEmpty(t *testing.T) {
	rb := NewRingBuffer[int](3)
	if !rb.IsEmpty() {
		t.Fatal("new buffer should be empty")
	}

	rb.Push(1)
	if rb.IsEmpty() {
		t.Fatal("buffer with items should not be empty")
	}
}

func TestRingBuffer_Clear(t *testing.T) {
	rb := NewRingBuffer[int](3)
	rb.Push(1)
	rb.Push(2)
	rb.Clear()

	if !rb.IsEmpty() {
		t.Fatal("buffer should be empty after Clear")
	}
	if rb.Cap() != 3 {
		t.Fatal("Cap should be unchanged after Clear")
	}
}

func TestRingBuffer_ToSlice(t *testing.T) {
	rb := NewRingBuffer[int](5)
	rb.Push(10)
	rb.Push(20)
	rb.Push(30)

	got := rb.ToSlice()
	want := []int{10, 20, 30}
	if len(got) != len(want) {
		t.Fatalf("ToSlice len = %d; want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ToSlice[%d] = %d; want %d", i, got[i], want[i])
		}
	}
}

func TestRingBuffer_ToSliceAfterOverflow(t *testing.T) {
	rb := NewRingBuffer[int](3)
	rb.Push(1)
	rb.Push(2)
	rb.Push(3)
	rb.Push(4) // evicts 1
	rb.Push(5) // evicts 2

	got := rb.ToSlice()
	want := []int{3, 4, 5}
	if len(got) != len(want) {
		t.Fatalf("ToSlice len = %d; want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ToSlice[%d] = %d; want %d", i, got[i], want[i])
		}
	}
}

func TestRingBuffer_ToSliceEmpty(t *testing.T) {
	rb := NewRingBuffer[int](3)
	got := rb.ToSlice()
	if len(got) != 0 {
		t.Fatalf("ToSlice(empty) = %v; want empty", got)
	}
}

func TestRingBuffer_ZeroCapacity(t *testing.T) {
	rb := NewRingBuffer[int](0)
	// Should default to capacity 1.
	if rb.Cap() != 1 {
		t.Fatalf("zero capacity Cap = %d; want 1", rb.Cap())
	}
}

func TestRingBuffer_WrapAround(t *testing.T) {
	rb := NewRingBuffer[int](3)
	// Fill and drain multiple times to test wraparound.
	for round := 0; round < 5; round++ {
		rb.Push(round*10 + 1)
		rb.Push(round*10 + 2)
		rb.Push(round*10 + 3)

		for j := 1; j <= 3; j++ {
			got, ok := rb.Pop()
			want := round*10 + j
			if !ok || got != want {
				t.Fatalf("round %d Pop = %d, %v; want %d, true", round, got, ok, want)
			}
		}
	}
}

func TestRingBuffer_Concurrent(t *testing.T) {
	rb := NewRingBuffer[int](50)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			rb.Push(n)
		}(i)
	}

	wg.Wait()
	// Buffer capacity is 50, so only 50 items remain.
	if rb.Len() != 50 {
		t.Fatalf("concurrent Len = %d; want 50", rb.Len())
	}
}
