package collections

import (
	"sort"
	"sync"
	"testing"
)

func TestSet_AddContains(t *testing.T) {
	s := NewSet[int]()
	if !s.Add(1) {
		t.Fatal("first Add should return true")
	}
	if s.Add(1) {
		t.Fatal("duplicate Add should return false")
	}
	if !s.Contains(1) {
		t.Fatal("Contains should return true for added item")
	}
	if s.Contains(2) {
		t.Fatal("Contains should return false for absent item")
	}
}

func TestSet_Remove(t *testing.T) {
	s := NewSet[string]()
	s.Add("a")

	if !s.Remove("a") {
		t.Fatal("Remove of existing item should return true")
	}
	if s.Remove("a") {
		t.Fatal("Remove of absent item should return false")
	}
	if s.Contains("a") {
		t.Fatal("removed item should not be contained")
	}
}

func TestSet_Len(t *testing.T) {
	s := NewSet[int]()
	if s.Len() != 0 {
		t.Fatalf("empty Len = %d; want 0", s.Len())
	}

	s.Add(1)
	s.Add(2)
	s.Add(2) // duplicate
	if s.Len() != 2 {
		t.Fatalf("Len = %d; want 2", s.Len())
	}
}

func TestSet_IsEmpty(t *testing.T) {
	s := NewSet[int]()
	if !s.IsEmpty() {
		t.Fatal("new set should be empty")
	}

	s.Add(1)
	if s.IsEmpty() {
		t.Fatal("set with items should not be empty")
	}
}

func TestSet_NewSetFrom(t *testing.T) {
	s := NewSetFrom(1, 2, 3, 2, 1)
	if s.Len() != 3 {
		t.Fatalf("NewSetFrom Len = %d; want 3", s.Len())
	}
	for _, v := range []int{1, 2, 3} {
		if !s.Contains(v) {
			t.Fatalf("should contain %d", v)
		}
	}
}

func TestSet_Values(t *testing.T) {
	s := NewSetFrom("c", "a", "b")
	vals := s.Values()
	sort.Strings(vals)
	if len(vals) != 3 || vals[0] != "a" || vals[1] != "b" || vals[2] != "c" {
		t.Fatalf("Values = %v; want [a b c]", vals)
	}
}

func TestSet_Clear(t *testing.T) {
	s := NewSetFrom(1, 2, 3)
	s.Clear()
	if !s.IsEmpty() {
		t.Fatal("set should be empty after Clear")
	}
}

func TestSet_Union(t *testing.T) {
	a := NewSetFrom(1, 2, 3)
	b := NewSetFrom(3, 4, 5)
	u := a.Union(b)

	if u.Len() != 5 {
		t.Fatalf("Union Len = %d; want 5", u.Len())
	}
	for _, v := range []int{1, 2, 3, 4, 5} {
		if !u.Contains(v) {
			t.Fatalf("Union should contain %d", v)
		}
	}
}

func TestSet_Intersection(t *testing.T) {
	a := NewSetFrom(1, 2, 3, 4)
	b := NewSetFrom(3, 4, 5, 6)
	inter := a.Intersection(b)

	if inter.Len() != 2 {
		t.Fatalf("Intersection Len = %d; want 2", inter.Len())
	}
	if !inter.Contains(3) || !inter.Contains(4) {
		t.Fatal("Intersection should contain 3 and 4")
	}
}

func TestSet_IntersectionEmpty(t *testing.T) {
	a := NewSetFrom(1, 2)
	b := NewSetFrom(3, 4)
	inter := a.Intersection(b)
	if inter.Len() != 0 {
		t.Fatalf("disjoint Intersection Len = %d; want 0", inter.Len())
	}
}

func TestSet_Difference(t *testing.T) {
	a := NewSetFrom(1, 2, 3, 4)
	b := NewSetFrom(3, 4, 5)
	diff := a.Difference(b)

	if diff.Len() != 2 {
		t.Fatalf("Difference Len = %d; want 2", diff.Len())
	}
	if !diff.Contains(1) || !diff.Contains(2) {
		t.Fatal("Difference should contain 1 and 2")
	}
	if diff.Contains(3) || diff.Contains(4) {
		t.Fatal("Difference should not contain shared elements")
	}
}

func TestSet_Range(t *testing.T) {
	s := NewSetFrom(1, 2, 3, 4, 5)
	count := 0
	s.Range(func(_ int) bool {
		count++
		return count < 3
	})
	if count != 3 {
		t.Fatalf("Range stopped after %d; want 3", count)
	}
}

func TestSet_RangeAll(t *testing.T) {
	s := NewSetFrom(10, 20, 30)
	var seen []int
	s.Range(func(v int) bool {
		seen = append(seen, v)
		return true
	})
	sort.Ints(seen)
	if len(seen) != 3 || seen[0] != 10 || seen[1] != 20 || seen[2] != 30 {
		t.Fatalf("Range saw %v; want [10 20 30]", seen)
	}
}

func TestSet_Concurrent(t *testing.T) {
	s := NewSet[int]()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			s.Add(n)
			s.Contains(n)
		}(i)
	}

	wg.Wait()
	if s.Len() != 100 {
		t.Fatalf("concurrent Len = %d; want 100", s.Len())
	}
}

func TestSet_UnionDoesNotMutate(t *testing.T) {
	a := NewSetFrom(1, 2)
	b := NewSetFrom(3, 4)
	a.Union(b)

	if a.Len() != 2 || b.Len() != 2 {
		t.Fatal("Union should not mutate operands")
	}
}

func TestSet_IntersectionSmallBig(t *testing.T) {
	// Ensure intersection works when 'other' is smaller.
	big := NewSetFrom(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	small := NewSetFrom(5, 10)
	inter := big.Intersection(small)
	if inter.Len() != 2 {
		t.Fatalf("Intersection Len = %d; want 2", inter.Len())
	}
}
