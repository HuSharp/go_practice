package geecache

import (
	"geecache/lru"
	"sync"
)

// cache 用来实例化 lru
type cache struct {
	mu 			sync.Mutex
	lru			*lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, val ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		// 判断了 c.lru 是否为 nil，如果等于 nil 再创建实例。
		// 这种方法称之为延迟初始化(Lazy Initialization)，
		// 一个对象的延迟初始化意味着该对象的创建将会延迟至第一次使用该对象时。主要用于提高性能，并减少程序内存要求。
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, val)
}

func (c *cache) get(key string) (val ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

