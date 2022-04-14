package synthetic

import (
	"strconv"
	"sync/atomic"
)

type counterGenerator struct {
	i int64
}

func (g *counterGenerator) Str() string {
	return strconv.FormatInt(atomic.AddInt64(&g.i, 1), 10)
}

func (g *counterGenerator) Int() int {
	return int(atomic.AddInt64(&g.i, 1))
}

// Counter returns a Generator giving a sequence of unique integers.
func Counter(start int) Generator {
	return &counterGenerator{
		i: int64(start),
	}
}
