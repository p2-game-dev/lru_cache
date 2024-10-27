package main

import (
	"container/list"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type CacheKey interface{}
type CacheValue interface{}

type LruCache struct {
	mutex               *sync.RWMutex
	capacity            int
	cache               map[CacheKey]CacheValue
	order               *list.List
	onRemoveKeyCallback func(key CacheKey)
}

func New(capacity int, onRemoveKeyCallback ...func(key CacheKey)) *LruCache {
	lr := &LruCache{
		capacity: capacity,
		cache:    make(map[CacheKey]CacheValue, capacity),
		order:    list.New(),
		mutex:    &sync.RWMutex{},
	}
	if len(onRemoveKeyCallback) > 0 {
		lr.onRemoveKeyCallback = onRemoveKeyCallback[0]
	}
	return lr
}

func (lru *LruCache) Set(key CacheKey, value CacheValue) error {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	_, exist := lru.cache[key]
	if exist {
		lru.cache[key] = value
		for el := lru.order.Front(); el != nil; el = el.Next() {
			if el.Value == key {
				lru.order.MoveToFront(el)
				break
			}
		}
		return nil
	}
	if len(lru.cache) >= lru.capacity {
		leastUsedKey := lru.order.Back().Value.(CacheValue)
		lru.order.Remove(lru.order.Back())
		delete(lru.cache, leastUsedKey)
		if lru.onRemoveKeyCallback != nil {
			lru.onRemoveKeyCallback(leastUsedKey)
		}
	}
	lru.cache[key] = value
	lru.order.PushFront(key)
	return nil
}

func (lru *LruCache) Get(key CacheKey) (CacheValue, error) {
	lru.mutex.RLock()
	defer lru.mutex.RUnlock()
	value, err := lru.cache[key]
	if !err {
		return nil, errors.New("key not found")
	}
	for el := lru.order.Front(); el != nil; el = el.Next() {
		if el.Value == key {
			lru.order.MoveToFront(el)
			break
		}
	}
	return value, nil
}

func (lru *LruCache) String() string {
	var builder strings.Builder
	for el := lru.order.Front(); el != nil; el = el.Next() {
		builder.WriteString(fmt.Sprintf("%v ", lru.cache[el.Value]))
	}
	return builder.String()
}

func (lru *LruCache) Close() {
	lru.onRemoveKeyCallback = nil
}
