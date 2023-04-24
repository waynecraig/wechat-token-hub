package cache

import "time"

type cacheItem struct {
	Type       string
	Value      string
	Expiration time.Time
}

var cache map[string]*cacheItem

func init() {
	cache = make(map[string]*cacheItem)
}

func GetCacheItem(itemType string) string {
	// check if the cache item is expired
	if cache[itemType] != nil && cache[itemType].Expiration.After(time.Now()) {
		// return the cache item
		return cache[itemType].Value
	}

	// if the cache item is expired or not set, return an empty string
	return ""
}

func SaveCacheItem(itemType string, value string, expiresIn int) {
	if expiresIn <= 0 {
		return
	}

	// calculate the expiration time
	expiration := time.Now().Add(time.Duration(expiresIn) * time.Second)

	// set the cache item to the new value and its expiration time
	cache[itemType] = &cacheItem{
		Type:       itemType,
		Value:      value,
		Expiration: expiration,
	}
}
