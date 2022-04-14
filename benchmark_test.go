package mkeeper

import (
	"strconv"
	"testing"
	"time"

	"github.com/svartvalp/mkeeper/synthetic"
)

const (
	testMaxSize = 512
)

type sameGenerator int

func (g sameGenerator) Str() string {
	return strconv.Itoa(int(g))
}

func (g sameGenerator) Int() int {
	return int(g)
}

func BenchmarkSame(b *testing.B) {
	g := sameGenerator(1)
	benchmarkCache(b, g)
}

func BenchmarkUniform(b *testing.B) {
	distintKeys := testMaxSize * 2
	g := synthetic.Uniform(0, distintKeys)
	benchmarkCache(b, g)
}

func BenchmarkUniformLess(b *testing.B) {
	distintKeys := testMaxSize
	g := synthetic.Uniform(0, distintKeys)
	benchmarkCache(b, g)
}

func BenchmarkCounter(b *testing.B) {
	g := synthetic.Counter(0)
	benchmarkCache(b, g)
}

func BenchmarkExponential(b *testing.B) {
	g := synthetic.Exponential(1.0)
	benchmarkCache(b, g)
}

func BenchmarkZipf(b *testing.B) {
	items := testMaxSize * 10
	g := synthetic.Zipf(0, items, 1.01)
	benchmarkCache(b, g)
}

func BenchmarkHotspot(b *testing.B) {
	items := testMaxSize * 2
	g := synthetic.Hotspot(0, items, 0.25)
	benchmarkCache(b, g)
}

func benchmarkCache(b *testing.B, g synthetic.Generator) {
	c := NewCache(WithMaxCapacity(testMaxSize), WithTTL(100*time.Millisecond), WithShardsCount(256))

	intCh := make(chan int, 2048)
	go func(n int) {
		for i := 0; i < n; i++ {
			intCh <- g.Int()
		}
	}(b.N)
	defer close(intCh)

	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(32)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := <-intCh
			_, ok := c.GetIfPresent(k)
			if !ok {
				c.Put(k, k)
			}
		}
	})
}
