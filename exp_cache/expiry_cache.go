package exp_cache

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	mut  *sync.RWMutex
	data map[uint64]*entry

	expL *list.List
	ttl  int64
}

func NewExpCache(ttl int64) *Cache {
	return &Cache{
		mut:  &sync.RWMutex{},
		data: make(map[uint64]*entry),
		expL: list.New(),
		ttl:  ttl,
	}
}

func (c *Cache) Put(key interface{}, val interface{}, hash uint64) {
	c.mut.Lock()
	defer c.mut.Unlock()
	en := c.data[hash]
	if en != nil {
		en.Key = key
		en.Value = val
		en.Exp = time.Now().UnixNano() + c.ttl
		if en.ExpEl != nil {
			c.expL.MoveToBack(en.ExpEl)
		} else {
			en.ExpEl = c.expL.PushBack(en)
		}
		return
	}
	en = &entry{
		Key:   key,
		Value: val,
		Hash:  hash,
		Exp:   time.Now().UnixNano() + c.ttl,
		ExpEl: nil,
	}
	c.data[hash] = en
	en.ExpEl = c.expL.PushBack(en)
}

func (c *Cache) Get(key interface{}, hash uint64) (interface{}, bool) {
	c.mut.RLock()
	defer c.mut.RUnlock()
	en := c.data[hash]
	if en == nil {
		return nil, false
	}

	if c.ttl > 0 && en.Exp <= time.Now().UnixNano() {
		return nil, false
	}

	if en.Key != key {
		return nil, false
	}

	return en.Value, true
}

func (c *Cache) Delete(key interface{}, hash uint64) {
	c.mut.Lock()
	defer c.mut.Unlock()
	en := c.data[hash]
	if en.Key != key {
		return
	}
	if en != nil && en.ExpEl != nil {
		c.expL.Remove(en.ExpEl)
	}
	delete(c.data, hash)
}

func (c *Cache) DeleteByHash(hash uint64) {
	c.mut.Lock()
	defer c.mut.Unlock()
	en := c.data[hash]
	if en != nil && en.ExpEl != nil {
		c.expL.Remove(en.ExpEl)
	}
	delete(c.data, hash)
}

func (c *Cache) Cleanup() []uint64 {
	cleaned := make([]uint64, 0)
	now := time.Now().UnixNano()
	c.mut.Lock()
	defer c.mut.Unlock()
	el := c.expL.Front()
	for el != nil {
		en := el.Value.(*entry)
		if en.Exp <= now {
			next := el.Next()
			saved := c.data[en.Hash]
			if saved == en {
				delete(c.data, en.Hash)
			}
			c.expL.Remove(el)
			el = next
			cleaned = append(cleaned, en.Hash)
		} else {
			return cleaned
		}
	}

	return cleaned
}
