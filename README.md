# go-collections

Type-safe concurrent data structures for Go, built with generics (Go 1.22+).

## Data Structures

| Type | Description |
|------|-------------|
| `SyncMap[K, V]` | Thread-safe map backed by `sync.Map` with full type safety |
| `Queue[T]` | Unbounded or bounded FIFO queue |
| `Set[T]` | Hash set with union, intersection, and difference operations |
| `RingBuffer[T]` | Fixed-size circular buffer that overwrites oldest entries when full |
| `Pool[T]` | Object pool with configurable factory and max size |

## Install

```bash
go get github.com/CorvidLabs/go-collections
```

## Usage

```go
import "github.com/CorvidLabs/go-collections"

// Concurrent map
m := collections.NewSyncMap[string, int]()
m.Store("count", 42)
v, ok := m.Load("count") // 42, true

// Bounded queue
q := collections.NewBoundedQueue[string](100)
q.Push("task-1")
item, ok := q.Pop() // "task-1", true

// Set with operations
a := collections.NewSetFrom(1, 2, 3)
b := collections.NewSetFrom(3, 4, 5)
union := a.Union(b)         // {1, 2, 3, 4, 5}
inter := a.Intersection(b)  // {3}

// Ring buffer (fixed-size, overwrites oldest)
rb := collections.NewRingBuffer[string](3)
rb.Push("a")
rb.Push("b")
rb.Push("c")
rb.Push("d") // evicts "a"

// Object pool
pool := collections.NewPool(func() []byte {
    return make([]byte, 4096)
}, 20)
buf := pool.Get()
defer pool.Put(buf)
```

## Thread Safety

All types are safe for concurrent use. Synchronization is handled internally via `sync.Mutex`, `sync.RWMutex`, or `sync.Map`.

## Build

```bash
make build   # compile
make test    # run tests with race detector
make bench   # run benchmarks
make lint    # vet + format check
```

## License

MIT
