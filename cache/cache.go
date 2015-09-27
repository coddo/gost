package cache

import (
	"bytes"
	"net/url"
	"time"
)

const (
	STATUS_ON  = true
	STATUS_OFF = false
)

const (
	CACHE_EXPIRE_TIME = 1 * time.Minute
)

var Status bool = STATUS_OFF
var selectedCacheExpireTime time.Duration

type Cacher interface {
	Cache()
	Invalidate()
	InvalidateIfExpired(limit time.Time)
	ResetExpireTime()
}

type Cache struct {
	Query       string
	Data        []byte
	StatusCode  int
	ContentType string
	File        string
	ExpireTime  time.Time
}

func (cache *Cache) Cache() {
	cache.ResetExpireTime()
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

func (cache *Cache) ResetExpireTime() {
	cache.ExpireTime = time.Now().Add(selectedCacheExpireTime)
}

func QueryByKey(key string) *Cache {
	go func() {
		getKeyChannel <- key
	}()

	flag := make(chan *Cache)
	defer close(flag)

	getChan <- flag

	return <-flag
}

func QueryByRequest(form url.Values, endpoint string) *Cache {
	return QueryByKey(MapKey(form, endpoint))
}

func MapKey(form url.Values, endpoint string) string {
	var buf bytes.Buffer

	buf.WriteString(endpoint)
	buf.WriteRune(':')
	buf.WriteString(form.Encode())

	return buf.String()
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
			break Loop
		case key := <-invalidateChan:
			invalidate(key)
		case cache := <-cacheChan:
			storeOrUpdate(cache)
		case flag := <-getChan:
			key := <-getKeyChannel

			item := memoryCache[key]
			if item != nil {
				item.ResetExpireTime()
			}

			flag <- item
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
	selectedCacheExpireTime = cacheExpireTime

	go startCachingLoop()
	go startExpiredInvalidator(cacheExpireTime)

	Status = STATUS_ON
}

func StopCachingSystem() {
	stopCachingSystem()

	go func() {
		exitChan <- 1
		close(exitChan)
	}()
}
