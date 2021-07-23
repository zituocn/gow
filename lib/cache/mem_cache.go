/*
一个内存缓存操作的 封装
sam
 */

package cache

import (
	"fmt"
	"reflect"
	"time"

	cache "github.com/patrickmn/go-cache"
)

// MemCache MemCache
type MemCache struct {
	cc *cache.Cache
}

const (
	defaultHour     = 72
	cleanupInterval = 1
)

// NewMemCache NewMemCache
func NewMemCache() *MemCache {
	mc := &MemCache{
		cc: cache.New(defaultHour*time.Hour, cleanupInterval*time.Minute),
	}
	return mc
}

// GetAll 获取obj的所有缓存数据
func (m *MemCache) GetAll(obj ICache) (data interface{}, err error) {
	key := obj.KeyName()        //获取key
	data, isExist := m.get(key) //取值
	//数据不存在
	if !isExist {
		data, err = obj.GetAllData()
		if err != nil {
			return
		}
		m.set(key, data)
	}
	return
}

// GetItem 取[]data中的一项数据
func (m *MemCache) GetItem(obj ICache, id int64) (item interface{}, err error) {
	data, err := m.GetAll(obj)
	if err != nil {
		return
	}
	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Slice {
		v := reflect.ValueOf(data)
		idStr := fmt.Sprintf("%v", id)
		length := v.Len()
		for i := 0; i < length; i++ {
			itemData := v.Index(i).Interface()
			if obj.PrimaryKey(itemData) == idStr {
				item = itemData
				break
			}
		}
	} else {
		err = fmt.Errorf(fmt.Sprintf("%#v 未缓存slice数据", obj))
	}
	return
}

// Remove 移出key对应的缓存
func (m *MemCache) Remove(obj ICache) {
	key := obj.KeyName()
	m.delete(key)
}

///===================private func =========================

//G et get data by string key
func (m *MemCache) get(key string) (data interface{}, isExist bool) {
	if key == "" {
		return
	}
	data, isExist = m.cc.Get(key)
	return
}

// Set Set
func (m *MemCache) set(key string, data interface{}) {
	if key == "" {
		return
	}
	m.cc.Set(key, data, cache.DefaultExpiration)
}

// clearAll ClearAll
func (m *MemCache) clearAll() {
	m.cc.Flush()
}

//delete 删除
func (m *MemCache) delete(key string) {
	if key == "" {
		return
	}
	m.cc.Delete(key)
}
