package cache

import (
	"bytes"
	"net/url"
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

func QueryCache(key string) *Cache {
	flag := make(chan int)
	defer func() {
		getChan <- flag
	}()

	getChan <- flag
	<-flag

	return memoryCache[key]
}

var memoryCache = make(map[string]*Cache)

var (
	getChan        = make(chan chan int)
	cacheChan      = make(chan *Cache)
	invalidateChan = make(chan string)
	exitChan       = make(chan int)
)

var exited bool = false

func stopCachingSystem() {
	exited = true

	close(getChan)
	close(cacheChan)
	close(invalidateChan)
	close(exitChan)
}

func invalidate(key string) {
	delete(memoryCache, key)
}

func storeOrUpdate(cache *Cache) {
	memoryCache[cache.Query] = cache
}

func startCachingLoop() {
loop:
	for {
		select {
		case <-exitChan:
			stopCachingSystem()
			break loop
		case key := <-invalidateChan:
			invalidate(key)
		case cache := <-cacheChan:
			storeOrUpdate(cache)
		case flag := <-getChan:
			flag <- 1
			<-getChan
		}
	}
}

func startExpiredInvalidator(cacheExpireTime time.Duration) {
	for !exited {
		time.Sleep(cacheExpireTime)

		if !exited {
			m := memoryCache
			date := time.Now()

			for _, item := range m {
				item.InvalidateIfExpired(date)
			}
		}
	}
}

func StartCachingSystem(cacheExpireTime time.Duration) {
	go startCachingLoop()
	go startExpiredInvalidator(cacheExpireTime)
}

func MapKey(form url.Values, method string, endpoint string) string {
	var buf bytes.Buffer

	buf.WriteString(endpoint)
	buf.WriteRune(':')
	buf.WriteString(method)
	buf.WriteRune('-')
	buf.WriteString(form.Encode())

	return buf.String()
}
