package sketch

import (
	"github.com/svartvalp/mkeeper/util"
)

const depth = 4

type Sketch struct {
	counters []uint64
	mask     uint32
}

func NewMinSketch(width int) *Sketch {
	c := &Sketch{}
	size := util.NextPowerOfTwo(uint32(width)) >> 2
	if size < 1 {
		size = 1
	}
	c.mask = size - 1
	c.counters = make([]uint64, size)
	return c
}

func (c *Sketch) Add(h uint64) {
	h1, h2 := uint32(h), uint32(h>>32)

	for i := uint32(0); i < depth; i++ {
		idx, off := c.position(h1 + i*h2)
		c.increment(idx, (16*i)+off)
	}
}

func (c *Sketch) Estimate(h uint64) uint8 {
	h1, h2 := uint32(h), uint32(h>>32)

	var min uint8 = 0xFF
	for i := uint32(0); i < depth; i++ {
		idx, off := c.position(h1 + i*h2)
		count := c.count(idx, (16*i)+off)
		if count < min {
			min = count
		}
	}
	return min
}

func (c *Sketch) Reset() {
	for i, v := range c.counters {
		if v != 0 {
			c.counters[i] = (v >> 1) & 0x7777777777777777
		}
	}
}

func (c *Sketch) position(h uint32) (idx uint32, off uint32) {
	idx = (h >> 2) & c.mask
	off = (h & 3) << 2
	return
}

func (c *Sketch) increment(idx, off uint32) {
	v := c.counters[idx]
	count := uint8(v>>off) & 0x0F
	if count < 15 {
		c.counters[idx] = v + (1 << off)
	}
}

func (c *Sketch) count(idx, off uint32) uint8 {
	v := c.counters[idx]
	return uint8(v>>off) & 0x0F
}
