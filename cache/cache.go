package cache

import (
	"time"
)

const (
	CACHE_EXPIRE_TIME = 1 * time.Minute
)

type Cacher interface {
	Cache()
	Invalidate()
	InvalidateIfExpired(limit time.Time)
}

type Cache struct {
	Query      string
	Data       []byte
	ExpireTime time.Time
}

func (cache *Cache) Cache() {
	cacheChan <- cache
}

func (cache *Cache) Invalidate() {
	invalidateChan <- cache.Query
}

func (cache *Cache) InvalidateIfExpired(limit time.Time) {
	if cache.ExpireTime.Before(limit) {
		cache.Invalidate()
	}
}

func QueryCache(query string) *Cache {
	flag := make(chan int)
	defer func() {
		getChan <- flag
	}()

	getChan <- flag
	<-flag

	return memoryCache[query]
}

var memoryCache = make(map[string]*Cache)

var (
	getChan        = make(chan chan int)
	cacheChan      = make(chan *Cache)
	invalidateChan = make(chan string)
	exitChan       = make(chan int)
)

var exited bool = false

func startCachingLoop() {
loop:
	for {
		select {
		case <-exitChan:
			exited = true

			close(getChan)
			close(cacheChan)
			close(invalidateChan)
			close(exitChan)

			break loop
		case query := <-invalidateChan:
			delete(memoryCache, query)
		case cache := <-cacheChan:
			memoryCache[cache.Query] = cache
		case flag := <-getChan:
			flag <- 1
			<-getChan
		}
	}
}

func startExpiredInvalidator() {
	for !exited {
		time.Sleep(time.Minute * 30)

		if !exited {
			m := memoryCache
			date := time.Now()

			for _, item := range m {
				item.InvalidateIfExpired(date)
			}
		}
	}
}

func StartCachingSystem() {
	go startCachingLoop()
	go startExpiredInvalidator()
}
