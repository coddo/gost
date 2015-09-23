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
	go func() {
		getKeyChannel <- key
	}()

	flag := make(chan *Cache)
	defer close(flag)

	getChan <- flag

	return <-flag
}

var memoryCache = make(map[string]*Cache)

var (
	getKeyChannel  = make(chan string)
	getChan        = make(chan chan *Cache)
	cacheChan      = make(chan *Cache)
	invalidateChan = make(chan string)
	exitChan       = make(chan int)
)

var exited bool = false

func stopCachingSystem() {
	exited = true

	close(getKeyChannel)
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
Loop:
	for {
		select {
		case <-exitChan:
			stopCachingSystem()
			break Loop
		case key := <-invalidateChan:
			invalidate(key)
		case cache := <-cacheChan:
			storeOrUpdate(cache)
		case flag := <-getChan:
			key := <-getKeyChannel
			flag <- memoryCache[key]
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

func StopCachingSystem() {
	go func() {
		exitChan <- 1
	}()
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
