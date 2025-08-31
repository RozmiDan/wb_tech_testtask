package lru_cache

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

// helper: must equal
func mustEqual[T comparable](t *testing.T, got, want T, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %v, want %v", msg, got, want)
	}
}

// helper: must be <=
func mustLE(t *testing.T, got, max int, msg string) {
	t.Helper()
	if got > max {
		t.Fatalf("%s: got %d, max %d", msg, got, max)
	}
}

func TestNewPanicsOnBadCapacity(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on zero/negative capacity")
		}
	}()
	_ = NewLruCache[string, int](0, 0)
}

func TestPutGetBasic(t *testing.T) {
	t.Parallel()

	c := NewLruCache[string, int](3, -1)
	c.Put("a", 1)
	c.Put("b", 2)

	mustEqual(t, c.Size(), 2, "size after 2 puts")

	mustEqual(t, c.Get("a"), 1, "get a")
	mustEqual(t, c.Get("b"), 2, "get b")
	mustEqual(t, c.Get("missing"), -1, "get missing returns default")
}

func TestUpdateMovesToFrontAndKeepsSize(t *testing.T) {
	t.Parallel()

	c := NewLruCache[string, int](2, -1)
	c.Put("a", 1)
	c.Put("b", 2)

	// update existing key "a"
	c.Put("a", 42)
	mustEqual(t, c.Size(), 2, "size unchanged after update")
	mustEqual(t, c.Get("a"), 42, "updated value visible")

	// Now insert "c" -> should evict LRU = "b"
	c.Put("c", 3)
	mustEqual(t, c.Size(), 2, "size after eviction")
	mustEqual(t, c.Get("a"), 42, "a survives")
	mustEqual(t, c.Get("b"), -1, "b evicted")
	mustEqual(t, c.Get("c"), 3, "c present")
}

func TestEvictionPolicyLRU(t *testing.T) {
	t.Parallel()

	c := NewLruCache[string, int](2, -1)
	// Insert a, b  (MRU=b, LRU=a)
	c.Put("a", 1)
	c.Put("b", 2)

	// Touch a to make it MRU -> order: MRU=a, LRU=b
	_ = c.Get("a")

	// Put c -> evict LRU=b
	c.Put("c", 3)

	mustEqual(t, c.Get("b"), -1, "b should be evicted")
	mustEqual(t, c.Get("a"), 1, "a should remain")
	mustEqual(t, c.Get("c"), 3, "c present")
}

func TestAllIterationOrder(t *testing.T) {
	t.Parallel()

	c := NewLruCache[string, int](3, -1)
	// After these puts: MRU=c, then b, then a (LRU)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)

	// Touch a -> MRU=a, then c, then b
	_ = c.Get("a")

	var keys []string
	for k, _v := range c.All() {
		_ = _v // value not used here
		keys = append(keys, k)
	}
	// Expect MRU-first order: a, c, b
	want := []string{"a", "c", "b"}
	if len(keys) != len(want) {
		t.Fatalf("iteration length mismatch: got %v, want %v", keys, want)
	}
	for i := range want {
		if keys[i] != want[i] {
			t.Fatalf("iteration order: got %v, want %v", keys, want)
		}
	}
}

func TestConcurrentPutGet(t *testing.T) {
	capacity := 128
	c := NewLruCache[int, int](capacity, -1)

	nG := runtime.GOMAXPROCS(0) * 4
	var wg sync.WaitGroup
	wg.Add(nG)

	opsPerG := 10_000

	for g := 0; g < nG; g++ {
		go func(seed int64) {
			defer wg.Done()
			r := rand.New(rand.NewSource(seed))
			for i := 0; i < opsPerG; i++ {
				k := r.Intn(capacity * 4) // ключи шире capacity → будет вытеснение
				if r.Intn(100) < 50 {
					c.Put(k, k*k)
				} else {
					_ = c.Get(k)
				}

				if i%2000 == 0 {
					time.Sleep(time.Microsecond)
				}
			}
		}(int64(g + 1))
	}

	wg.Wait()

	// После гонки проверим несколько инвариантов
	mustLE(t, c.Size(), capacity, "size must not exceed capacity")

	// Быстрая проверка, что кэш «живой»: случайные ключи не паникнут
	for k := 0; k < 10; k++ {
		_ = c.Get(k)
	}
}

func TestConcurrentAllVsWrites(t *testing.T) {
	// Проверим, что All() с RLock не падает при параллельных записях
	c := NewLruCache[int, int](64, -1)

	// Писатель
	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		for {
			select {
			case <-stop:
				return
			default:
				c.Put(i%128, i)
				i++
			}
		}
	}()

	// Читатель итератором
	iterAttempts := 200
	for a := 0; a < iterAttempts; a++ {
		count := 0
		for k, v := range c.All() {
			_ = k
			_ = v
			count++
			if count > 1000 { // не зависаем на слишком длинной итерации
				break
			}
		}
	}

	close(stop)
	wg.Wait()

	mustLE(t, c.Size(), 64, "size must not exceed capacity after concurrent use")
}
