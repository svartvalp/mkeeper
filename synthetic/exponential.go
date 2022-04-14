package synthetic

import (
	"math/rand"
	"strconv"
	"time"
)

type exponentialGenerator struct {
	r    *rand.Rand
	mean float64
}

func (g *exponentialGenerator) Str() string {
	return strconv.Itoa(g.Int())
}

func (g *exponentialGenerator) Int() int {
	return int(g.r.ExpFloat64() * g.mean)
}

// Exponential returns a Generator resembling an exponential distribution.
func Exponential(mean float64) Generator {
	return &exponentialGenerator{
		r:    rand.New(rand.NewSource(time.Now().UnixNano())),
		mean: mean,
	}
}
