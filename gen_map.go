package main

import (
	"sync"
)

// simple generic sync map for concurrent lookups

type GenMap[K comparable, V any] struct {
	speciesMap sync.Map
}

func (sm *GenMap[K, V]) Store(key K, value V) {
	sm.speciesMap.Store(key, value)
}

func (sm *GenMap[K, V]) Load(key K) (V, bool) {
	var val V
	if value, ok := sm.speciesMap.Load(key); ok {
		return value.(V), ok
	}
	return val, false

}

func (sm *GenMap[K, V]) Range(f func(key K, value V) bool) {
	sm.speciesMap.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}
