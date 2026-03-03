package collections

import (
	"sync"
	"testing"
)

func TestPool_GetPut(t *testing.T) {
	created := 0
	p := NewPool(func() int {
		created++
		return created
	}, 10)

	// First Get creates via factory.
	v := p.Get()
	if v != 1 {
		t.Fatalf("first Get = %d; want 1", v)
	}

	// Put it back and Get again.
	p.Put(v)
	v = p.Get()
	if v != 1 {
		t.Fatalf("Get after Put = %d; want 1 (reused)", v)
	}

	// Pool is empty, next Get creates new.
	v = p.Get()
	if v != 2 {
		t.Fatalf("Get from empty pool = %d; want 2 (new)", v)
	}
}

func TestPool_MaxSize(t *testing.T) {
	p := NewPool(func() int { return 0 }, 2)

	if !p.Put(1) {
		t.Fatal("first Put should succeed")
	}
	if !p.Put(2) {
		t.Fatal("second Put should succeed")
	}
	if p.Put(3) {
		t.Fatal("third Put should fail (at capacity)")
	}
}

func TestPool_Unlimited(t *testing.T) {
	p := NewPool(func() int { return 0 }, 0)

	for i := 0; i < 100; i++ {
		if !p.Put(i) {
			t.Fatalf("Put(%d) failed on unlimited pool", i)
		}
	}

	if p.Len() != 100 {
		t.Fatalf("Len = %d; want 100", p.Len())
	}
}

func TestPool_Len(t *testing.T) {
	p := NewPool(func() int { return 0 }, 10)
	if p.Len() != 0 {
		t.Fatalf("empty Len = %d; want 0", p.Len())
	}

	p.Put(1)
	p.Put(2)
	if p.Len() != 2 {
		t.Fatalf("Len = %d; want 2", p.Len())
	}

	p.Get()
	if p.Len() != 1 {
		t.Fatalf("Len after Get = %d; want 1", p.Len())
	}
}

func TestPool_Clear(t *testing.T) {
	p := NewPool(func() int { return 0 }, 10)
	p.Put(1)
	p.Put(2)
	p.Clear()

	if p.Len() != 0 {
		t.Fatalf("Len after Clear = %d; want 0", p.Len())
	}
}

func TestPool_Prefill(t *testing.T) {
	counter := 0
	p := NewPool(func() int {
		counter++
		return counter
	}, 5)

	p.Prefill(3)
	if p.Len() != 3 {
		t.Fatalf("Len after Prefill(3) = %d; want 3", p.Len())
	}

	// Prefill respects maxSize.
	p.Prefill(10)
	if p.Len() != 5 {
		t.Fatalf("Len after Prefill(10) = %d; want 5 (capped)", p.Len())
	}
}

func TestPool_PrefillUnlimited(t *testing.T) {
	p := NewPool(func() int { return 42 }, 0)
	p.Prefill(50)
	if p.Len() != 50 {
		t.Fatalf("unlimited Prefill Len = %d; want 50", p.Len())
	}
}

func TestPool_NegativeMaxSize(t *testing.T) {
	p := NewPool(func() int { return 0 }, -5)
	// Should treat as unlimited.
	for i := 0; i < 10; i++ {
		if !p.Put(i) {
			t.Fatalf("Put(%d) failed; negative maxSize should be unlimited", i)
		}
	}
}

func TestPool_Concurrent(t *testing.T) {
	p := NewPool(func() []byte {
		return make([]byte, 1024)
	}, 50)

	p.Prefill(20)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := p.Get()
			// Simulate work.
			buf[0] = 1
			p.Put(buf)
		}()
	}

	wg.Wait()
	// Pool should have items returned (up to maxSize).
	if p.Len() == 0 {
		t.Fatal("pool should have items after concurrent use")
	}
	if p.Len() > 50 {
		t.Fatalf("pool should not exceed maxSize; got %d", p.Len())
	}
}

func TestPool_FactoryCalled(t *testing.T) {
	calls := 0
	p := NewPool(func() string {
		calls++
		return "item"
	}, 5)

	// Each Get on empty pool calls factory.
	p.Get()
	p.Get()
	p.Get()

	if calls != 3 {
		t.Fatalf("factory called %d times; want 3", calls)
	}
}

func TestPool_LIFO(t *testing.T) {
	p := NewPool(func() int { return 0 }, 10)
	p.Put(1)
	p.Put(2)
	p.Put(3)

	// Pool returns last-in first (stack behavior for better cache locality).
	v := p.Get()
	if v != 3 {
		t.Fatalf("Get = %d; want 3 (LIFO)", v)
	}
}
