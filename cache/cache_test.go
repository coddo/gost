package cache

import (
	"encoding/json"
	"gost/config"
	"testing"
	"time"
)

const TESTS_CACHE_EXPIRE_TIME = 5 * time.Second

const (
	TESTS_QUERY1        = "test:x%2==0"
	TESTS_QUERY2        = "test:(x+y)%z>1"
	TESTS_QUERY3        = "test:z>550"
	TESTS_QUERY_EXPIRED = "thisNeedsToExpire"
)

type CacheTest struct {
	X int
	Y int
	Z int
}

func TestCache(t *testing.T) {
	var items []CacheTest
	var cachedItems []*Cache
	var expiringItem *Cache

	config.InitTestsDatabase()
	StartCachingSystem(TESTS_CACHE_EXPIRE_TIME)
	defer StopCachingSystem()

	items = testInitItems(t)

	testFetchInexistentCache(t, TESTS_QUERY1)
	cachedItems, expiringItem = testAddingToCache(t, items)
	testFetchingFromCache(t, cachedItems)
	testRemovingFromCache(t, cachedItems)
	testFetchInexistentCache(t, TESTS_QUERY2)

	time.Sleep(TESTS_CACHE_EXPIRE_TIME)
	testExpiringItem(t, expiringItem)
}

func testExpiringItem(t *testing.T, expiringItem *Cache) {
	it := QueryCache(expiringItem.Query)

	if it != nil {
		t.Fatal("The cache items did not properly expire")
	}
}

func testFetchInexistentCache(t *testing.T, mockQuery string) {
	// Will never be added
	data := QueryCache("keySFAFSAGKAGHAJSKfhaskfhaskf")
	if data != nil {
		t.Fatal("Unexpected output from cache")
	}

	// Will be added later during the test
	data = QueryCache(mockQuery)
	if data != nil {
		t.Fatal("Unexpected output from cache")
	}
}

func testFetchingFromCache(t *testing.T, cachedItems []*Cache) {
	var q1 *Cache
	var q2 *Cache
	var q3 *Cache
	i := 0

	for i < 2 {
		q1 = QueryCache(cachedItems[0].Query)
		q2 = QueryCache(cachedItems[1].Query)
		q3 = QueryCache(cachedItems[2].Query)

		if q1 == nil || q2 == nil || q3 == nil {
			t.Fatal("Cache didn't properly return test items")
		}

		i++
	}

	if q1.Query != cachedItems[0].Query || q2.Query != cachedItems[1].Query || q3.Query != cachedItems[2].Query {
		t.Fatal("Wrong cache values were returned")
	}
}

func testAddingToCache(t *testing.T, items []CacheTest) ([]*Cache, *Cache) {
	var cachedItems []*Cache
	var expiringCacheItem *Cache

	q1 := make([]CacheTest, 0)
	q2 := make([]CacheTest, 0)
	q3 := make([]CacheTest, 0)

	expireTime := time.Now().Add(TESTS_CACHE_EXPIRE_TIME)

	// First type
	for i := 0; i < len(items); i++ {
		if items[i].X%2 == 0 {
			q1 = append(q1, items[i])
		}
	}
	j1, _ := json.MarshalIndent(q1, "", "  ")
	c1 := &Cache{
		Query:      TESTS_QUERY1,
		Data:       j1,
		ExpireTime: expireTime,
	}
	c1.Cache()
	cachedItems = append(cachedItems, c1)

	// Second type
	for i := 0; i < len(items); i++ {
		if ((items[i].X + items[i].Y) % items[i].Z) > 1 {
			q2 = append(q2, items[i])
		}
	}
	j2, _ := json.MarshalIndent(q2, "", "  ")
	c2 := &Cache{
		Query:      TESTS_QUERY2,
		Data:       j2,
		ExpireTime: expireTime,
	}
	c2.Cache()
	cachedItems = append(cachedItems, c2)

	// Third type
	for i := 0; i < len(items); i++ {
		if items[i].Z > 550 {
			q3 = append(q3, items[i])
		}
	}
	j3, _ := json.MarshalIndent(q3, "", "  ")
	c3 := &Cache{
		Query:      TESTS_QUERY3,
		Data:       j3,
		ExpireTime: expireTime,
	}
	c3.Cache()
	cachedItems = append(cachedItems, c3)

	// Expiring type
	expiringCacheItem = &Cache{
		Query:      TESTS_QUERY_EXPIRED,
		Data:       j1,
		ExpireTime: expireTime,
	}
	expiringCacheItem.Cache()

	return cachedItems, expiringCacheItem
}

func testRemovingFromCache(t *testing.T, cachedItems []*Cache) {
	for _, it := range cachedItems {
		it.Invalidate()
	}
}

func testInitItems(t *testing.T) []CacheTest {
	var items []CacheTest

	for i := 1; i < 1000; i++ {
		items = append(items, CacheTest{
			X: i,
			Y: i * 11 / 3,
			Z: i * 3,
		})
	}

	return items
}
