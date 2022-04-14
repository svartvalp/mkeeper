package policy

import (
	"sync"

	"github.com/svartvalp/mkeeper"
	"github.com/svartvalp/mkeeper/bloom"
	"github.com/svartvalp/mkeeper/sketch"
)

type Policy struct {
	mu        *sync.RWMutex
	sketch    *sketch.Sketch
	bloom     *bloom.Filter
	samples   uint64
	threshold uint64
	keys      map[uint64]struct{}
}

func NewPolicy(cap int64) *Policy {
	return &Policy{
		mu:        &sync.RWMutex{},
		sketch:    sketch.NewMinSketch(int(cap)),
		bloom:     bloom.NewBloomFilter(int(2*cap), 0.1),
		samples:   0,
		threshold: uint64(8 * cap),
		keys:      make(map[uint64]struct{}),
	}
}

func (p *Policy) ProcessEvents(events []mkeeper.Event) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, e := range events {
		switch e.Type {
		case mkeeper.Delete:
			delete(p.keys, e.Hash)
		case mkeeper.Get:
			if p.bloom.Put(e.Hash) {
				p.sketch.Add(e.Hash)
			}
			p.samples++
		case mkeeper.Miss:
		default:
			if p.bloom.Put(e.Hash) {
				p.sketch.Add(e.Hash)
			}
			p.keys[e.Hash] = struct{}{}
			p.samples++
		}
	}
	if p.samples >= p.threshold {
		p.sketch.Reset()
		p.bloom.Reset()
		p.samples = 0
	}
}

func (p *Policy) Victim(h uint64) uint64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	cur := 0
	estimate := p.estimate(h)
	for k := range p.keys {
		victimEs := p.estimate(k)
		if victimEs < estimate {
			return k
		}
		cur++
		if cur > 5 {
			return h
		}
	}
	return h
}

func (p *Policy) estimate(h uint64) uint8 {
	freq := p.sketch.Estimate(h)
	if p.bloom.Contains(h) {
		freq++
	}
	return freq
}
