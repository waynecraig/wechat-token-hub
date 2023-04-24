package cache

import (
	"testing"
	"time"
)

func TestGetCacheItem(t *testing.T) {
	// set up a cache item
	cache["test"] = &cacheItem{
		Type:       "test",
		Value:      "value",
		Expiration: time.Now().Add(time.Duration(1) * time.Hour),
	}

	// test getting an existing cache item
	result := GetCacheItem("test")
	if result != "value" {
		t.Errorf("GetCacheItem returned %s, expected value", result)
	}

	// test getting a non-existent cache item
	result = GetCacheItem("nonexistent")
	if result != "" {
		t.Errorf("GetCacheItem returned %s, expected empty string", result)
	}

	// test getting an expired cache item
	cache["expired"] = &cacheItem{
		Type:       "expired",
		Value:      "value",
		Expiration: time.Now().Add(time.Duration(-1) * time.Hour),
	}
	result = GetCacheItem("expired")
	if result != "" {
		t.Errorf("GetCacheItem returned %s, expected empty string", result)
	}
}

func TestSaveCacheItem(t *testing.T) {
	// test saving a cache item
	SaveCacheItem("test", "value", 3600)
	if cache["test"] == nil {
		t.Errorf("SaveCacheItem did not save cache item")
	}

	// test saving a cache item with a negative expiration time
	SaveCacheItem("negative", "value", -3600)
	if cache["negative"] != nil {
		t.Errorf("SaveCacheItem saved cache item with negative expiration time")
	}
}
