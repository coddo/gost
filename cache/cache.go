package cache

import (
	"bytes"
	"errors"
	"gost/util"
	"time"
)

const (
	// StatusON shows that the cashing system is up and running
	StatusON = true
	// StatusOFF shows that the caching system is stopped and not functional
	StatusOFF = false
)

const (
	// CacheExpireTime represents the maximum duration that an item can stay cached
	CacheExpireTime = 7 * 24 * time.Hour
)

var (
	// ErrKeyInvalidated is used when a search key is inexistent or has been invalidated
	ErrKeyInvalidated = errors.New("The search key has been invalidated")

	// ErrKeyFormat is used when the format of the key is wrong or cannot be parsed
	ErrKeyFormat = errors.New("The search key is not in a correct format")

	// ErrCachingSystemStopped is used to indicate that the caching system is not available or stopped
	ErrCachingSystemStopped = errors.New("The search key has been invalidated")
)

var (
	// Status represents the current status of the caching system
	Status = StatusOFF

	selectedCacheExpireTime time.Duration
)

var memoryCache = make(map[string]*Cache)

var (
	getKeyChannel  = make(chan string)
	getChan        = make(chan chan *Cache)
	errorChan      = make(chan error)
	cacheChan      = make(chan *Cache)
	invalidateChan = make(chan string)
	exitChan       = make(chan int)
)

// Cacher is the interface that has all the basic methods used by cached items
//
// Each cache item can either: become Cached, become Invalidated (triggered or expired),
// or have its expire time reset
type Cacher interface {
	Cache()
	Invalidate()
	InvalidateIfExpired(limit time.Time)
	ResetExpireTime()
}

// A Cache entity is used to store precise information in the memory cache
// using a key (unique idenfier) and its actual data
type Cache struct {
	Key         string
	Data        []byte
	StatusCode  int
	ContentType string
	File        string
	ExpireTime  time.Time
}

// Cache caches the current entity into memory
func (cache *Cache) Cache() {
	go func() {
		cacheChan <- cache
	}()
}

// Invalidate invalidates the current entity by removing it from the cache
func (cache *Cache) Invalidate() {
	go func() {
		invalidateChan <- cache.Key
	}()
}

// InvalidateIfExpired checks whether the active time of the current entity
// has expired, and if it did, it invalidates it
func (cache *Cache) InvalidateIfExpired(limit time.Time) {
	if util.IsDateExpired(cache.ExpireTime, limit) {
		cache.Invalidate()
	}
}

// ResetExpireTime resets the timer for when the current entity will expire
func (cache *Cache) ResetExpireTime() {
	go func() {
		cache.ExpireTime = util.NextDateFromNow(selectedCacheExpireTime)
	}()
}

// StartCachingSystem starts the caching system. The status is set to StatusON.
//
// The following async loops are started:
// - Selection loop for the request channels
// - Invalidator loop which makes sure that the entities will not be stored in the cache forever
func StartCachingSystem(cacheExpireTime time.Duration) {
	Status = StatusON

	selectedCacheExpireTime = cacheExpireTime

	go startCachingLoop()
	go startExpiredInvalidator(cacheExpireTime)
}

// StopCachingSystem stops the caching system/ The status is set to StatusOFF.
// All the async loops are stopped and the channels (except the errors channel) are closed.
func StopCachingSystem() {
	Status = StatusOFF

	stopCachingSystem()

	go func() {
		exitChan <- 1
		close(exitChan)
	}()
}

// QueryByKey searches for a certain storage key in the memory cache
// and returns the found Cache item, or an error if it is inexistent
// or there was a problem with the search
func QueryByKey(key string) (*Cache, error) {
	if Status == StatusOFF {
		return nil, ErrCachingSystemStopped
	}

	go func() {
		getKeyChannel <- key
	}()

	flagChan := make(chan *Cache)
	defer close(flagChan)

	getChan <- flagChan

	select {
	case returnItem := <-flagChan:
		return returnItem, nil
	case err := <-errorChan:
		return nil, err
	}
}

// QueryByRequest maps a given endpoint to a storage key and then searches for
// that key in the memory cache and returns the found Cache item, or an error
// if it is inexistent or there was a problem with the search
func QueryByRequest(endpoint string) (*Cache, error) {
	return QueryByKey(MapKey(endpoint))
}

// MapKey creates a storage key out of a string, used for identifying cache items.
func MapKey(endpoint string) string {
	var buf bytes.Buffer

	buf.WriteString(endpoint)

	return buf.String()
}

func stopCachingSystem() {
	close(getKeyChannel)
	close(getChan)
	close(cacheChan)
	close(invalidateChan)
}

func invalidate(key string) {
	if _, exists := memoryCache[key]; exists {
		delete(memoryCache, key)
	}
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

			if len(key) == 0 {
				errorChan <- ErrKeyFormat
			}

			if item, ok := memoryCache[key]; ok {
				item.ResetExpireTime()
				flag <- item
			} else {
				errorChan <- ErrKeyInvalidated
			}
		}
	}
}

func startExpiredInvalidator(cacheExpireTime time.Duration) {
	for Status == StatusON {
		date := util.Now()

		for _, item := range memoryCache {
			item.InvalidateIfExpired(date)
		}

		time.Sleep(cacheExpireTime)
	}
}
