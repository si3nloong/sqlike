package lrucache

import (
	"reflect"
	"sync"
)

// Cache interface
type Cache[V any] interface {
	Get(reflect.Type) (V, bool)
	Set(reflect.Type, V)
	Size() int
	Len() int
	Purge()
}

type listItem[V any] struct {
	key   reflect.Type
	value V
	prev  *listItem[V]
	next  *listItem[V]
}

// FIXME: Supposingly we should make key as generic type,
// but current golang version (up to v1.18) doesn't support generic for comparable.
// See this : https://github.com/golang/go/issues/51179
type lru[V any] struct {
	// maximum entry of the lru cache
	size int

	// mutex lock is to prevent concurrency race
	mu sync.Mutex

	// store the current entry of the list
	current *listItem[V]

	// indexes for the lru cache entries, for fast searching purpose
	appendix map[reflect.Type]*listItem[V]

	// linked list for the lru cache entries
	list *listItem[V]
}

// Returns a lru cache with a size which comply to the `Cache` interface.
//
//   lrucache.New[string](100)
func New[V any](capacity int) Cache[V] {
	if capacity <= 1 {
		panic("lrucache size should be at least 2")
	}
	return &lru[V]{
		appendix: make(map[reflect.Type]*listItem[V]),
		size:     capacity,
	}
}

// Returns the value using key.
//
//   lrucache.Get(reflect.TypeOf(time.Time{}))
func (l *lru[V]) Get(t reflect.Type) (V, bool) {
	v, ok := l.appendix[t]
	if ok {
		return v.value, true
	}
	return *new(V), false
}

// Set the key and value into cache.
//
//   lrucache.Set(reflect.TypeOf(time.Time{}), time.Now())
func (l *lru[V]) Set(t reflect.Type, value V) {
	li := new(listItem[V])
	li.key = t
	li.value = value
	l.mu.Lock()
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
	l.mu.Unlock()
}

// Len returns the number of items in the cache.
func (l *lru[V]) Len() int {
	return len(l.appendix)
}

// Size returns the size of the cache.
func (l *lru[V]) Size() int {
	return l.size
}

// Purge clear all the entries inside the cache.
func (l *lru[V]) Purge() {
	l.mu.Lock()
	l.current = nil
	l.list = nil
	for k := range l.appendix {
		delete(l.appendix, k)
	}
	l.mu.Unlock()
}
