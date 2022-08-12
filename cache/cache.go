// Package cache impl the cache for gocq
package cache

import (
	"os"
	"sync"

	"github.com/AkiraXie/go-cqhttp-btree-manager/btree"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Cache wraps the btree.DB for concurrent safe
type Cache struct {
	lock sync.RWMutex
	db   *btree.DB
}

// Insert 添加媒体缓存
func (c *Cache) Insert(md5, data []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var hash [16]byte
	copy(hash[:], md5)
	c.db.Insert(&hash[0], data)
}

// Get 获取缓存信息
func (c *Cache) Get(md5 []byte) []byte {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var hash [16]byte
	copy(hash[:], md5)
	return c.db.Get(&hash[0])
}

// Delete 删除指定缓存
func (c *Cache) Delete(md5 []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var hash [16]byte
	copy(hash[:], md5)
	_ = c.db.Delete(&hash[0])
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || errors.Is(err, os.ErrExist)
}

// Init 初始化 Cache
func Init(file string) (cache *Cache) {
	cache = &Cache{}
	if PathExists(file) {
		db, err := btree.Open(file)
		if err != nil {
			log.Fatalf("open cache failed: %v", err)
		}
		cache.db = db
	} else {
		db, err := btree.Create(file)
		if err != nil {
			log.Fatalf("create cache failed: %v", err)
		}
		cache.db = db
	}
	return
}

func (c *Cache) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.db.Close()
}

func (c *Cache) Foreach(iter func(key [16]byte, value []byte)) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	c.db.Foreach(iter)
}
