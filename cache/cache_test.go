package cache

import (
	"encoding/json"
	testconfig "gost/tests/config"
	"gost/util"
	"testing"
	"time"
)

type cacheTest struct {
	X int
	Y int
	Z int
}

func TestCache(t *testing.T) {
	const cacheExpireTime = time.Millisecond * 700

	var cacheKeys = []combinedKey{
		combinedKey{Key: "/testKey1", DataKey: "testDataKey1"},
		combinedKey{Key: "/testKey2", DataKey: "testDataKey2"},
		combinedKey{Key: "/testTheThirdKey", DataKey: "testTheThirdDataKey"},
		combinedKey{Key: "/thisNeedsToExpire", DataKey: "thisNeedsToDataExpire"},
	}

	var items []cacheTest
	var cachedItems []*Cache
	var expiringItem *Cache

	testconfig.InitTestsDatabase()
	StartCachingSystem(cacheExpireTime)
	defer StopCachingSystem()

	items = testInitItems(t)

	testFetchInexistentCache(t, cacheKeys[0])
	cachedItems, expiringItem = testAddingToCache(t, items, cacheKeys)
	testFetchingFromCache(t, cachedItems)
	testRemovingFromCache(t, cachedItems)
	testFetchInexistentCache(t, cacheKeys[1])
	testExpiringItem(t, expiringItem, cacheExpireTime)
}

func testExpiringItem(t *testing.T, expiringItem *Cache, cacheExpireTime time.Duration) {
	t.Log("[info] Testing the expired cache invalidation system")

	time.Sleep(cacheExpireTime * 2)

	_, err := Query(expiringItem.Key, expiringItem.DataKey)

	if err == nil || err != ErrKeyInvalidated {
		t.Fatal("The cache items did not properly expire")
	}
}

func testFetchInexistentCache(t *testing.T, mockQuery combinedKey) {
	t.Log("[info] Testing the cache querying system with inexistent or invalid data")

	// Will never be added
	var inexistentKey, _ = util.GenerateUUID()
	var inexistentDataKey, _ = util.GenerateUUID()
	data, _ := Query(inexistentKey, inexistentDataKey)
	if data != nil {
		t.Fatal("[error] Unexpected output from cache")
	}

	// Will be added later during the test
	data, _ = Query(mockQuery.Key, mockQuery.DataKey)
	if data != nil {
		t.Fatal("[error] Unexpected output from cache")
	}
}

func testFetchingFromCache(t *testing.T, cachedItems []*Cache) {
	t.Log("[info] Testing the cache querying system with valid data")

	var q1 *Cache
	var q2 *Cache
	var q3 *Cache
	i := 0

	for i < 2 {
		q1, _ = Query(cachedItems[0].Key, cachedItems[0].DataKey)
		q2, _ = Query(cachedItems[1].Key, cachedItems[1].DataKey)
		q3, _ = Query(cachedItems[2].Key, cachedItems[2].DataKey)

		if q1 == nil || q2 == nil || q3 == nil {
			t.Fatal("[error] Cache didn't properly return test items")
		}

		i++
	}

	if q1.Key != cachedItems[0].Key || q2.Key != cachedItems[1].Key || q3.Key != cachedItems[2].Key {
		t.Fatal("[error] Wrong cache values were returned")
	}
}

func testAddingToCache(t *testing.T, items []cacheTest, cacheKeys []combinedKey) ([]*Cache, *Cache) {
	t.Log("[info] Testing the data caching system")

	var cachedItems = make([]*Cache, 3)
	var expiringCacheItem *Cache

	var q1 []cacheTest
	var q2 []cacheTest
	var q3 []cacheTest

	// First type
	for i := 0; i < len(items); i++ {
		if items[i].X%2 == 0 {
			q1 = append(q1, items[i])
		}
	}
	j1, _ := json.MarshalIndent(q1, "", "  ")
	c1 := &Cache{
		Key:     cacheKeys[0].Key,
		DataKey: cacheKeys[0].DataKey,
		Data:    j1,
	}
	cachedItems[0] = c1

	// Second type
	for i := 0; i < len(items); i++ {
		if ((items[i].X + items[i].Y) % items[i].Z) > 1 {
			q2 = append(q2, items[i])
		}
	}
	j2, _ := json.MarshalIndent(q2, "", "  ")
	c2 := &Cache{
		Key:     cacheKeys[1].Key,
		DataKey: cacheKeys[1].DataKey,
		Data:    j2,
	}
	cachedItems[1] = c2

	// Third type
	for i := 0; i < len(items); i++ {
		if items[i].Z > 550 {
			q3 = append(q3, items[i])
		}
	}
	j3, _ := json.MarshalIndent(q3, "", "  ")
	c3 := &Cache{
		Key:     cacheKeys[2].Key,
		DataKey: cacheKeys[2].DataKey,
		Data:    j3,
	}
	cachedItems[2] = c3

	// Expiring type
	expiringCacheItem = &Cache{
		Key:     cacheKeys[3].Key,
		DataKey: cacheKeys[3].DataKey,
		Data:    j1,
	}

	expiringCacheItem.Cache()
	for _, cachedItem := range cachedItems {
		cachedItem.Cache()
	}

	time.Sleep(500 * time.Millisecond)

	return cachedItems, expiringCacheItem
}

func testRemovingFromCache(t *testing.T, cachedItems []*Cache) {
	t.Log("[info] Testing the cache invalidation system")

	for _, it := range cachedItems {
		it.Invalidate()
	}
}

func testInitItems(t *testing.T) []cacheTest {
	var items []cacheTest

	for i := 1; i < 1000; i++ {
		items = append(items, cacheTest{
			X: i,
			Y: i * 11 / 3,
			Z: i * 3,
		})
	}

	return items
}
