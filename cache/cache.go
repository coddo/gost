package cache

import (
	"bytes"
	"errors"
	"log"
	"time"
)

const (
	STATUS_ON  = true
	STATUS_OFF = false
)

const (
	CACHE_EXPIRE_TIME = 7 * 24 * time.Hour
)

var (
	KEY_INVALIDATED_ERROR = errors.New("The search key has been invalidated")
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
	Key         string
	Data        []byte
	StatusCode  int
	ContentType string
	File        string
	ExpireTime  time.Time
}

func (cache *Cache) Cache() {
	go func() {
		cacheChan <- cache
	}()
}

func (cache *Cache) Invalidate() {
	go func() {
		invalidateChan <- cache.Key
	}()
}

func (cache *Cache) InvalidateIfExpired(limit time.Time) {
	if cache.ExpireTime.Before(limit) {
		cache.Invalidate()
	}
}

func (cache *Cache) ResetExpireTime() {
	go func() {
		cache.ExpireTime = time.Now().Add(selectedCacheExpireTime)
	}()
}

func QueryByKey(key string) *Cache {
	go func() {
		getKeyChannel <- key
	}()

	flagChan := make(chan *Cache)
	defer close(flagChan)

	getChan <- flagChan

	log.Println("*****GET CHANNEL SET CORRECTLY")
	select {
	case returnItem := <-flagChan:
		log.Println("*****ITEM SUCCESSFULLY RETRIEVED FROM CHANNEL: ", returnItem)
		return returnItem
	case <-errorChan:
		log.Println("*****ITEM RETRIEVE FAILED FROM CHANNEL")
		return nil
	}
}

func QueryByRequest(endpoint string) *Cache {
	return QueryByKey(MapKey(endpoint))
}

func MapKey(endpoint string) string {
	var buf bytes.Buffer

	buf.WriteString(endpoint)

	return buf.String()
}

var memoryCache = make(map[string]*Cache)

var (
	getKeyChannel  = make(chan string)
	getChan        = make(chan chan *Cache)
	errorChan      = make(chan error)
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
	close(errorChan)
}

func invalidate(key string) {
	delete(memoryCache, key)
}

func storeOrUpdate(cache *Cache) {
	cache.ResetExpireTime()
	memoryCache[cache.Key] = cache
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

			if item, ok := memoryCache[key]; ok {
				log.Println("*****MAP SEARCH WAS OK")
				item.ResetExpireTime()
				flag <- item
			} else {
				log.Println("*****MAP SEARCH WAS A FAILURE")
				errorChan <- KEY_INVALIDATED_ERROR
			}
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
