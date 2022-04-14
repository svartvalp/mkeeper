package mkeeper

import (
	"sync/atomic"
	"time"

	"github.com/svartvalp/mkeeper/exp_cache"
	"github.com/svartvalp/mkeeper/hash"
	"github.com/svartvalp/mkeeper/policy"
	"github.com/svartvalp/mkeeper/util"
)

const workersCount = 4

type cache struct {
	shards      map[int64]*exp_cache.Cache
	shardsCount int64
	shardsMask  int64

	ttl  int64
	cap  int64
	size int64

	cleanCh   chan int64
	cleanTick *time.Ticker

	listener *EventListener
	policy   *policy.Policy

	closeCh chan struct{}
}

func NewCache(opts ...Option) *cache {
	cfg := defaultConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	shardsCount := int64(util.NextPowerOfTwo(cfg.ShardsCount))
	shardsMask := shardsCount - 1
	shards := make(map[int64]*exp_cache.Cache, shardsCount)
	for i := int64(0); i < shardsCount; i++ {
		shards[i] = exp_cache.NewExpCache(int64(cfg.TTL))
	}

	cleanBuf := shardsCount
	if cfg.CleanBuf > 0 {
		cleanBuf = cfg.CleanBuf
	}

	pol := policy.NewPolicy(cfg.MaxCap)
	listener := NewEventListener(
		[]EventProcessor{cfg.Stats, pol},
		cfg.EventBuf,
		cfg.EventTick,
	)
	listener.Start()
	c := &cache{
		shards:      shards,
		shardsCount: shardsCount,
		shardsMask:  shardsMask,
		ttl:         int64(cfg.TTL),
		cap:         cfg.MaxCap,
		size:        0,
		cleanCh:     make(chan int64, cleanBuf),
		cleanTick:   time.NewTicker(cfg.CleanTick),
		listener:    listener,
		policy:      pol,
	}
	c.runCleaner()
	return c
}

func defaultConfig() *Config {
	return &Config{
		ShardsCount: 256,
		CleanTick:   time.Second,
		EventBuf:    256,
		EventTick:   100 * time.Millisecond,
		Stats: &Stats{
			OnPut:    func() {},
			OnDelete: func() {},
			OnGet:    func() {},
			OnMiss:   func() {},
		},
	}
}

type Option func(c *Config)

type Config struct {
	ShardsCount uint32
	TTL         time.Duration
	MaxCap      int64
	CleanTick   time.Duration
	CleanBuf    int64
	EventBuf    int
	EventTick   time.Duration
	Stats       *Stats
}

func (c *cache) GetIfPresent(key interface{}) (interface{}, bool) {
	h := hash.H(key)
	shard := c.shard(h)

	res, ok := shard.Get(key, h)
	c.listener.Send(Event{
		Type: Get,
		Hash: h,
	})
	if !ok {
		c.listener.Send(Event{
			Type: Miss,
			Hash: h,
		})
	}
	return res, ok
}

func (c *cache) Put(key interface{}, value interface{}) {
	h := hash.H(key)
	if c.cap > 0 && c.size >= c.cap {
		victim := c.policy.Victim(h)
		if victim == h {
			return
		}
		c.deleteByHash(victim)
	}

	shard := c.shard(h)

	shard.Put(key, value, h)
	atomic.AddInt64(&c.size, 1)
	c.listener.Send(Event{
		Type: Put,
		Hash: h,
	})
}

func (c *cache) Invalidate(key interface{}) {
	h := hash.H(key)
	shard := c.shard(h)
	shard.Delete(key, h)
	atomic.AddInt64(&c.size, -1)
	c.listener.Send(Event{
		Type: Delete,
		Hash: h,
	})
}

func (c *cache) deleteByHash(h uint64) {
	shard := c.shard(h)
	shard.DeleteByHash(h)
	atomic.AddInt64(&c.size, -1)
	c.listener.Send(Event{
		Type: Delete,
		Hash: h,
	})
}

func (c *cache) shard(h uint64) *exp_cache.Cache {
	return c.shards[int64(h)&c.shardsMask]
}

func (c *cache) runCleaner() {
	go func() {
		for {
			select {
			case <-c.cleanTick.C:
				for in := range c.shards {
					c.cleanCh <- in
				}
			case <-c.closeCh:
				close(c.cleanCh)
				return
			}
		}
	}()
	for i := 0; i < workersCount; i++ {
		go func() {
			for v := range c.cleanCh {
				cleaned := c.shards[v].Cleanup()
				if len(cleaned) > 0 {
					atomic.AddInt64(&c.size, int64(-len(cleaned)))
					for _, h := range cleaned {
						c.listener.Send(Event{
							Type: Delete,
							Hash: h,
						})
					}
					select {
					case c.cleanCh <- v:
					default:
					}
				}
			}
		}()
	}
}
