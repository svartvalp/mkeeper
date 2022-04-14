package bloom

import (
	"math"

	"github.com/svartvalp/mkeeper/util"
)

type Filter struct {
	numHashes uint32
	bitsMask  uint32
	bits      []uint64
}

func NewBloomFilter(ins int, fpp float64) *Filter {
	f := &Filter{}
	ln2 := math.Log(2.0)
	factor := -math.Log(fpp) / (ln2 * ln2)

	numBits := util.NextPowerOfTwo(uint32(float64(ins) * factor))
	if numBits == 0 {
		numBits = 1
	}
	f.bitsMask = numBits - 1

	if ins == 0 {
		f.numHashes = 1
	} else {
		f.numHashes = uint32(ln2 * float64(numBits) / float64(ins))
	}

	size := int(numBits+63) / 64
	f.bits = make([]uint64, size)
	return f
}

func (f *Filter) Put(h uint64) bool {
	h1, h2 := uint32(h), uint32(h>>32)
	var o uint = 1
	for i := uint32(0); i < f.numHashes; i++ {
		o &= f.set((h1 + (i * h2)) & f.bitsMask)
	}
	return o == 1
}

func (f *Filter) Contains(h uint64) bool {
	h1, h2 := uint32(h), uint32(h>>32)
	var o uint = 1
	for i := uint32(0); i < f.numHashes; i++ {
		o &= f.get((h1 + (i * h2)) & f.bitsMask)
	}
	return o == 1
}

func (f *Filter) set(i uint32) uint {
	idx, shift := i/64, i%64
	val := f.bits[idx]
	mask := uint64(1) << shift
	f.bits[idx] |= mask
	return uint((val & mask) >> shift)
}

func (f *Filter) get(i uint32) uint {
	idx, shift := i/64, i%64
	val := f.bits[idx]
	mask := uint64(1) << shift
	return uint((val & mask) >> shift)
}

func (f *Filter) Reset() {
	for i := range f.bits {
		f.bits[i] = 0
	}
}
