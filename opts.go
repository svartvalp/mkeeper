package mkeeper

import (
	"time"
)

func WithTTL(ttl time.Duration) Option {
	return func(c *Config) {
		c.TTL = ttl
	}
}

func WithMaxCapacity(cap int64) Option {
	return func(c *Config) {
		c.MaxCap = cap
	}
}

func WithCleanTick(tick time.Duration) Option {
	return func(c *Config) {
		c.CleanTick = tick
	}
}

func WithCleanBuf(buf int64) Option {
	return func(c *Config) {
		c.CleanBuf = buf
	}
}

func WithStats(stats *Stats) Option {
	return func(c *Config) {
		c.Stats = stats
	}
}

func WithShardsCount(count uint32) Option {
	return func(c *Config) {
		c.ShardsCount = count
	}
}
