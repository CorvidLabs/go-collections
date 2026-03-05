package collections

import (
	"sync"
	"testing"
)

func TestLRU_GetPut(t *testing.T) {
	c := NewLRU[string, int](3)

	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)

	for _, tc := range []struct {
		key  string
		want int
	}{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	} {
		got, ok := c.Get(tc.key)
		if !ok || got != tc.want {
			t.Fatalf("Get(%q) = %d, %v; want %d, true", tc.key, got, ok, tc.want)
		}
	}
}

func TestLRU_GetMiss(t *testing.T) {
	c := NewLRU[string, int](2)
	v, ok := c.Get("missing")
	if ok || v != 0 {
		t.Fatalf("Get(missing) = %d, %v; want 0, false", v, ok)
	}
}

func TestLRU_Eviction(t *testing.T) {
	c := NewLRU[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)

	// "a" is the LRU; inserting "c" should evict it.
	ek, ev, evicted := c.Put("c", 3)
	if !evicted || ek != "a" || ev != 1 {
		t.Fatalf("Put(c) evicted = %q:%d, %v; want a:1, true", ek, ev, evicted)
	}

	if c.Contains("a") {
		t.Fatal("evicted key 'a' should not be in cache")
	}

	if v, ok := c.Get("b"); !ok || v != 2 {
		t.Fatalf("Get(b) = %d, %v; want 2, true", v, ok)
	}
	if v, ok := c.Get("c"); !ok || v != 3 {
		t.Fatalf("Get(c) = %d, %v; want 3, true", v, ok)
	}
}

func TestLRU_AccessUpdatesRecency(t *testing.T) {
	c := NewLRU[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)

	// Access "a" to make it recently used.
	c.Get("a")

	// Now "b" is LRU; inserting "c" should evict "b".
	ek, _, evicted := c.Put("c", 3)
	if !evicted || ek != "b" {
		t.Fatalf("expected eviction of 'b', got %q (evicted=%v)", ek, evicted)
	}
}

func TestLRU_UpdateExisting(t *testing.T) {
	c := NewLRU[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)

	// Update "a" (should not evict anything).
	_, _, evicted := c.Put("a", 10)
	if evicted {
		t.Fatal("updating existing key should not evict")
	}

	v, ok := c.Get("a")
	if !ok || v != 10 {
		t.Fatalf("Get(a) = %d, %v; want 10, true", v, ok)
	}

	if c.Len() != 2 {
		t.Fatalf("Len = %d; want 2", c.Len())
	}
}

func TestLRU_Peek(t *testing.T) {
	c := NewLRU[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)

	// Peek should not update recency.
	v, ok := c.Peek("a")
	if !ok || v != 1 {
		t.Fatalf("Peek(a) = %d, %v; want 1, true", v, ok)
	}

	// "a" is still LRU; inserting "c" should evict "a".
	ek, _, _ := c.Put("c", 3)
	if ek != "a" {
		t.Fatalf("expected eviction of 'a' (Peek should not update recency), got %q", ek)
	}
}

func TestLRU_PeekMiss(t *testing.T) {
	c := NewLRU[string, int](2)
	v, ok := c.Peek("missing")
	if ok || v != 0 {
		t.Fatalf("Peek(missing) = %d, %v; want 0, false", v, ok)
	}
}

func TestLRU_Remove(t *testing.T) {
	c := NewLRU[string, int](3)
	c.Put("a", 1)
	c.Put("b", 2)

	v, ok := c.Remove("a")
	if !ok || v != 1 {
		t.Fatalf("Remove(a) = %d, %v; want 1, true", v, ok)
	}
	if c.Len() != 1 {
		t.Fatalf("Len = %d; want 1", c.Len())
	}
	if c.Contains("a") {
		t.Fatal("removed key should not be in cache")
	}
}

func TestLRU_RemoveMiss(t *testing.T) {
	c := NewLRU[string, int](2)
	v, ok := c.Remove("missing")
	if ok || v != 0 {
		t.Fatalf("Remove(missing) = %d, %v; want 0, false", v, ok)
	}
}

func TestLRU_Contains(t *testing.T) {
	c := NewLRU[string, int](2)
	c.Put("a", 1)

	if !c.Contains("a") {
		t.Fatal("Contains(a) should be true")
	}
	if c.Contains("b") {
		t.Fatal("Contains(b) should be false")
	}
}

func TestLRU_LenCap(t *testing.T) {
	c := NewLRU[int, int](5)
	if c.Cap() != 5 {
		t.Fatalf("Cap = %d; want 5", c.Cap())
	}
	if c.Len() != 0 {
		t.Fatalf("Len = %d; want 0", c.Len())
	}

	c.Put(1, 10)
	c.Put(2, 20)
	if c.Len() != 2 {
		t.Fatalf("Len = %d; want 2", c.Len())
	}
}

func TestLRU_Keys(t *testing.T) {
	c := NewLRU[string, int](3)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)

	// Most recent first: c, b, a
	keys := c.Keys()
	if len(keys) != 3 || keys[0] != "c" || keys[1] != "b" || keys[2] != "a" {
		t.Fatalf("Keys = %v; want [c b a]", keys)
	}

	// Access "a" to move it to front.
	c.Get("a")
	keys = c.Keys()
	if keys[0] != "a" || keys[1] != "c" || keys[2] != "b" {
		t.Fatalf("Keys after Get(a) = %v; want [a c b]", keys)
	}
}

func TestLRU_Clear(t *testing.T) {
	c := NewLRU[string, int](3)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Clear()

	if c.Len() != 0 {
		t.Fatalf("Len after Clear = %d; want 0", c.Len())
	}
	if c.Contains("a") || c.Contains("b") {
		t.Fatal("cache should be empty after Clear")
	}
}

func TestLRU_ZeroCapacity(t *testing.T) {
	c := NewLRU[int, int](0)
	if c.Cap() != 1 {
		t.Fatalf("Cap = %d; want 1 (minimum)", c.Cap())
	}
}

func TestLRU_SingleCapacity(t *testing.T) {
	c := NewLRU[string, int](1)
	c.Put("a", 1)

	ek, ev, evicted := c.Put("b", 2)
	if !evicted || ek != "a" || ev != 1 {
		t.Fatalf("eviction = %q:%d, %v; want a:1, true", ek, ev, evicted)
	}
	if c.Len() != 1 {
		t.Fatalf("Len = %d; want 1", c.Len())
	}
}

func TestLRU_Concurrent(t *testing.T) {
	c := NewLRU[int, int](100)
	var wg sync.WaitGroup

	// Concurrent writes.
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			c.Put(n, n*10)
		}(i)
	}

	// Concurrent reads.
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			c.Get(n)
		}(i)
	}

	wg.Wait()

	if c.Len() > 100 {
		t.Fatalf("Len = %d; should not exceed capacity 100", c.Len())
	}
}

func TestLRU_EvictionOrder(t *testing.T) {
	c := NewLRU[int, int](3)
	c.Put(1, 10)
	c.Put(2, 20)
	c.Put(3, 30)

	// Access order: 1 (oldest), then 2, then 3 (newest).
	// Access 1 and 2 to make 3 the LRU.
	c.Get(1)
	c.Get(2)

	// Now order is: 2 (most recent), 1, 3 (LRU).
	ek, _, _ := c.Put(4, 40)
	if ek != 3 {
		t.Fatalf("expected eviction of 3 (LRU), got %d", ek)
	}
}
