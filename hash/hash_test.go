package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testhash struct {
	A string
	B int
	C int64
}

func TestHash(t *testing.T) {
	require.Equal(t, uint64(3), H(3))
	require.NotPanicsf(t, func() {
		H(testhash{
			A: "1",
			B: 2,
			C: 3,
		})
	}, "panic on struct hash")
	require.NotPanicsf(t, func() {
		H(&testhash{
			A: "1",
			B: 2,
			C: 3,
		})
	}, "panic on pointer hash")
	hash1 := H(testhash{
		A: "1",
		B: 2,
		C: 3,
	})
	hash2 := H(testhash{
		A: "1",
		B: 2,
		C: 3,
	})
	require.Equal(t, hash1, hash2)
	hash1 = H(&testhash{
		A: "1",
		B: 2,
		C: 3,
	})
	hash2 = H(&testhash{
		A: "1",
		B: 2,
		C: 3,
	})
	require.Equal(t, hash1, hash2)
}
