package lrucache

import (
	"reflect"
	"sync"
)

type LRUCache[V any] interface {
	Get(reflect.Type) (V, bool)
	Set(reflect.Type, V) int
}

type lru[V any] struct {
	mu       sync.Mutex
	appendix map[reflect.Type]int
	stack    []V
}

func (l *lru[V]) Get(t reflect.Type) (V, bool) {
	v, ok := l.appendix[t]
	if ok {
		return l.stack[v], true
	}
	return *new(V), false
}

func (l *lru[V]) Set(t reflect.Type, value V) int {
	l.mu.Lock()
	l.stack = append(l.stack, value)
	pos := len(l.stack) - 1
	l.appendix[t] = pos
	l.mu.Unlock()
	return pos
}

func New[V any](capacity int) LRUCache[V] {
	return &lru[V]{
		appendix: make(map[reflect.Type]int),
		stack:    make([]V, capacity),
	}
}
