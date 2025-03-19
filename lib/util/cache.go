/*

memory cache

use GobEncode

*/

package util

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var (
	l sync.Mutex
)

// MemoryCache memory cache struct
type MemoryCache struct {
	defaultHour     int
	cleanupInterval int
	cc              *cache.Cache
}

// GobEncode gob encode
func GobEncode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	l.Lock()
	defer l.Unlock()
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("god encode:%v", err)
	}
	return buf.Bytes(), err
}

// GobDecode GobDecode
func GobDecode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	l.Lock()
	defer l.Unlock()
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

// NewMemoryCache return a new memory cache
func NewMemoryCache(defaultHour, cleanupInterval int) *MemoryCache {
	return &MemoryCache{
		defaultHour:     defaultHour,
		cleanupInterval: cleanupInterval,
		cc:              cache.New(time.Duration(defaultHour)*time.Hour, time.Duration(cleanupInterval)*time.Minute),
	}
}

// expireTime 缓存的有效时长，创建一个过期时间为expireTime的缓存
// cleanupInterval 清理过期缓存的频率，每cleanupInterval清理一次过期的缓存
func NewMemoryCacheWithExpireTime(expireTime, cleanupInterval time.Duration) *MemoryCache {
	return &MemoryCache{
		//defaultHour:     0,
		//cleanupInterval: 0,
		cc: cache.New(expireTime, cleanupInterval),
	}
}

// SetCache set cache to key and return error
func (m *MemoryCache) SetCache(key string, value interface{}) error {
	if key == "" {
		return errors.New("key is nil")
	}
	if value == nil {
		return errors.New("value is nil")
	}
	if m.cc == nil {
		return errors.New("must init memory cache first")
	}
	data, err := GobEncode(value)
	if err != nil {
		return err
	}
	return m.cc.Add(key, data, cache.DefaultExpiration)
}

// GetCache return cache data
func (m *MemoryCache) GetCache(key string, to interface{}) (ok bool, err error) {
	if key == "" {
		err = errors.New("key is nil")
		return
	}
	if m.cc == nil {
		err = errors.New("must init memory cache first")
		return
	}
	data, ok := m.cc.Get(key)
	if ok {
		err = GobDecode(data.([]byte), to)
		return
	}
	return
}

// RemoveCache delete cache
func (m *MemoryCache) RemoveCache(key string) error {
	if key == "" {
		return errors.New("key is nil")
	}
	if m.cc == nil {
		return errors.New("must init memory cache first")
	}
	m.cc.Delete(key)
	return nil
}
