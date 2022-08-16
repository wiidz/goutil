package memoryMng

import (
	"github.com/patrickmn/go-cache"
	"log"
	"time"
)

var cacheMng = cache.New(cache.DefaultExpiration, 0)

type MemoryMng struct {
	Client *cache.Cache
}

// NewCacheMng 获取系统缓存管理器
func NewCacheMng() *MemoryMng {
	var memoryMng = MemoryMng{
		Client: cacheMng,
	}
	return &memoryMng
}

// Get 提取
func (mng *MemoryMng) Get(keyName string) (data interface{}, isExist bool) {
	log.Println("get", keyName)
	data, isExist = mng.Client.Get(keyName)
	return
}

// GetString 提取
func (mng *MemoryMng) GetString(keyName string) (data string, isExist bool) {
	log.Println("GetString", keyName)
	var temp interface{}
	temp, isExist = mng.Client.Get(keyName)
	var ok bool
	if data, ok = temp.(string); !ok {
		return "", false
	}
	return
}

// Set 存储
func (mng *MemoryMng) Set(keyName string, data interface{}, expire time.Duration) {
	log.Println("set", keyName, data, expire)
	mng.Client.Set(keyName, data, expire)
}
