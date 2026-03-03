package collections

import (
	"sort"
	"sync"
	"testing"
)

func TestSyncMap_StoreAndLoad(t *testing.T) {
	m := NewSyncMap[string, int]()
	m.Store("a", 1)
	m.Store("b", 2)

	v, ok := m.Load("a")
	if !ok || v != 1 {
		t.Fatalf("Load(a) = %d, %v; want 1, true", v, ok)
	}

	v, ok = m.Load("b")
	if !ok || v != 2 {
		t.Fatalf("Load(b) = %d, %v; want 2, true", v, ok)
	}
}

func TestSyncMap_LoadMissing(t *testing.T) {
	m := NewSyncMap[string, int]()
	v, ok := m.Load("missing")
	if ok || v != 0 {
		t.Fatalf("Load(missing) = %d, %v; want 0, false", v, ok)
	}
}

func TestSyncMap_Overwrite(t *testing.T) {
	m := NewSyncMap[string, string]()
	m.Store("key", "old")
	m.Store("key", "new")

	v, ok := m.Load("key")
	if !ok || v != "new" {
		t.Fatalf("Load(key) = %q, %v; want 'new', true", v, ok)
	}
}

func TestSyncMap_LoadOrStore(t *testing.T) {
	m := NewSyncMap[string, int]()

	// First call stores.
	actual, loaded := m.LoadOrStore("x", 10)
	if loaded || actual != 10 {
		t.Fatalf("first LoadOrStore = %d, %v; want 10, false", actual, loaded)
	}

	// Second call loads existing.
	actual, loaded = m.LoadOrStore("x", 99)
	if !loaded || actual != 10 {
		t.Fatalf("second LoadOrStore = %d, %v; want 10, true", actual, loaded)
	}
}

func TestSyncMap_LoadAndDelete(t *testing.T) {
	m := NewSyncMap[string, int]()
	m.Store("k", 42)

	v, loaded := m.LoadAndDelete("k")
	if !loaded || v != 42 {
		t.Fatalf("LoadAndDelete = %d, %v; want 42, true", v, loaded)
	}

	_, ok := m.Load("k")
	if ok {
		t.Fatal("key should be deleted")
	}
}

func TestSyncMap_LoadAndDeleteMissing(t *testing.T) {
	m := NewSyncMap[int, string]()
	v, loaded := m.LoadAndDelete(1)
	if loaded || v != "" {
		t.Fatalf("LoadAndDelete(missing) = %q, %v; want '', false", v, loaded)
	}
}

func TestSyncMap_Delete(t *testing.T) {
	m := NewSyncMap[string, int]()
	m.Store("a", 1)
	m.Delete("a")

	_, ok := m.Load("a")
	if ok {
		t.Fatal("key should be deleted")
	}
}

func TestSyncMap_Range(t *testing.T) {
	m := NewSyncMap[int, string]()
	m.Store(1, "one")
	m.Store(2, "two")
	m.Store(3, "three")

	seen := make(map[int]string)
	m.Range(func(k int, v string) bool {
		seen[k] = v
		return true
	})

	if len(seen) != 3 {
		t.Fatalf("Range visited %d entries; want 3", len(seen))
	}
	if seen[1] != "one" || seen[2] != "two" || seen[3] != "three" {
		t.Fatalf("unexpected values: %v", seen)
	}
}

func TestSyncMap_RangeEarlyStop(t *testing.T) {
	m := NewSyncMap[int, int]()
	for i := 0; i < 10; i++ {
		m.Store(i, i)
	}

	count := 0
	m.Range(func(_ int, _ int) bool {
		count++
		return count < 3
	})

	if count != 3 {
		t.Fatalf("Range iterated %d times; want 3", count)
	}
}

func TestSyncMap_Len(t *testing.T) {
	m := NewSyncMap[string, int]()
	if m.Len() != 0 {
		t.Fatalf("empty Len = %d; want 0", m.Len())
	}

	m.Store("a", 1)
	m.Store("b", 2)
	if m.Len() != 2 {
		t.Fatalf("Len = %d; want 2", m.Len())
	}

	m.Delete("a")
	if m.Len() != 1 {
		t.Fatalf("Len after delete = %d; want 1", m.Len())
	}
}

func TestSyncMap_Keys(t *testing.T) {
	m := NewSyncMap[string, int]()
	m.Store("b", 2)
	m.Store("a", 1)
	m.Store("c", 3)

	keys := m.Keys()
	sort.Strings(keys)
	if len(keys) != 3 || keys[0] != "a" || keys[1] != "b" || keys[2] != "c" {
		t.Fatalf("Keys = %v; want [a b c]", keys)
	}
}

func TestSyncMap_Values(t *testing.T) {
	m := NewSyncMap[string, int]()
	m.Store("a", 1)
	m.Store("b", 2)

	values := m.Values()
	sort.Ints(values)
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Fatalf("Values = %v; want [1 2]", values)
	}
}

func TestSyncMap_Concurrent(t *testing.T) {
	m := NewSyncMap[int, int]()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			m.Store(n, n*2)
			m.Load(n)
		}(i)
	}

	wg.Wait()
	if m.Len() != 100 {
		t.Fatalf("concurrent Len = %d; want 100", m.Len())
	}
}

func TestSyncMap_EmptyKeys(t *testing.T) {
	m := NewSyncMap[string, string]()
	keys := m.Keys()
	if keys != nil {
		t.Fatalf("empty Keys = %v; want nil", keys)
	}
}

func TestSyncMap_EmptyValues(t *testing.T) {
	m := NewSyncMap[string, string]()
	values := m.Values()
	if values != nil {
		t.Fatalf("empty Values = %v; want nil", values)
	}
}
