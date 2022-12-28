package store

import (
	"log"
	"sync"
)

var store *redisStore
var once sync.Once
var mutex sync.Mutex

type Store interface {
	Get(key string) string
	Set(key string, value string)
}

type redisStore struct {
	store map[string]string
}

func (c redisStore) Get(key string) string {
	mutex.Lock()
	val := c.store[key]
	log.Printf("Value fetched: store[%s]", key)
	mutex.Unlock()
	return val
}

func (c redisStore) Set(key, value string) {
	mutex.Lock()
	c.store[key] = value
	log.Printf("Value stored: store[%s]=%s", key, value)
	mutex.Unlock()
}

func GetStore() Store {
	once.Do(func() {
		store = &redisStore{}
		store.store = make(map[string]string)
		log.Println("Store created")
	})
	return *store
}
