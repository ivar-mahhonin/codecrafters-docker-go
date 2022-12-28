package store

import (
	"log"
	"sync"
	"time"
)

var rStore *redisStore
var once sync.Once
var mutex sync.Mutex

type Store interface {
	Get(key string) (string, bool)
	Set(key string, value string, expiration time.Time)
}

type Value struct {
	value      string
	expiration time.Time
}
type redisStore struct {
	store map[string]Value
}

func (c redisStore) Get(key string) (string, bool) {
	mutex.Lock()
	defer mutex.Unlock()
	storeValue := c.store[key]
	log.Printf("Value fetched: store[%s]", key)
	if storeValue.isExpired() {
		delete(c.store, key)
		return storeValue.value, false
	}
	return storeValue.value, true
}

func (c redisStore) Set(key, value string, expiration time.Time) {
	mutex.Lock()
	defer mutex.Unlock()
	c.store[key] = Value{value, expiration}
	log.Printf("Value stored: store[%s]=%s, px=%s", key, value, expiration)
}

func GetStore() Store {
	once.Do(func() {
		rStore = &redisStore{}
		rStore.store = make(map[string]Value)
		log.Println("Store created")
	})
	return *rStore
}

func (v Value) isExpired() bool {
	if v.expiration.IsZero() {
		return false
	}
	return v.expiration.Before(time.Now())
}
