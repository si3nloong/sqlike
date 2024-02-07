package lrucache

import (
	"sync"
)

// Cache interface
type Cache[K comparable, V any] interface {
	Get(K) (V, bool)
	Set(K, V)
	Size() int
	Len() int
	Purge()
}

type listItem[K comparable, V any] struct {
	key   K
	value V
	prev  *listItem[K, V]
	next  *listItem[K, V]
}

type lru[K comparable, V any] struct {
	// maximum entry of the lru cache
	size int

	// mutex lock is to prevent concurrency race
	mu *sync.RWMutex

	// store the current entry of the list
	current *listItem[K, V]

	// indexes for the lru cache entries, for fast searching purpose
	appendix map[K]*listItem[K, V]

	// linked list for the lru cache entries
	list *listItem[K, V]
}

// Returns a lru cache with a size which comply to the `Cache` interface.
//
//	lrucache.New[string](100)
func New[K comparable, V any](capacity int) Cache[K, V] {
	if capacity <= 1 {
		panic("lrucache size should be at least 2")
	}
	return &lru[K, V]{
		mu:       new(sync.RWMutex),
		appendix: make(map[K]*listItem[K, V]),
		size:     capacity,
	}
}

// Returns the value using key.
//
//	lrucache.Get(reflect.TypeOf(time.Time{}))
func (l *lru[K, V]) Get(t K) (V, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	v, ok := l.appendix[t]
	if ok {
		return v.value, true
	}
	return *new(V), false
}

// Set the key and value into cache.
//
//	lrucache.Set(reflect.TypeOf(time.Time{}), time.Now())
func (l *lru[K, V]) Set(t K, value V) {
	li := new(listItem[K, V])
	li.key = t
	li.value = value
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list != nil {
		li.prev = l.list
		l.current.next = li
		l.current = li
	} else {
		// if cache have no items, then we set the `li` as the first entry
		l.current = li
		l.list = li
	}
	l.appendix[t] = li
	if len(l.appendix) > l.size {
		next := l.list.next
		delete(l.appendix, l.list.key)
		l.list = next
	}
}

// Len returns the number of items in the cache.
func (l *lru[K, V]) Len() int {
	return len(l.appendix)
}

// Size returns the size of the cache.
func (l *lru[K, V]) Size() int {
	return l.size
}

// Purge clear all the entries inside the cache.
func (l *lru[K, V]) Purge() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.current = nil
	l.list = nil
	for k := range l.appendix {
		delete(l.appendix, k)
	}
}
