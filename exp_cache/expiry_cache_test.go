package exp_cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPutGet(t *testing.T) {
	c := NewExpCache(int64(time.Minute))
	c.Put(1, 1, 1)
	c.Put(2, 2, 2)
	val, ok := c.Get(1, 1)
	require.Equal(t, true, ok)
	require.Equal(t, 1, val)
	val, ok = c.Get(2, 2)
	require.Equal(t, true, ok)
	require.Equal(t, 2, val)
	val, ok = c.Get(3, 3)
	require.Equal(t, false, ok)
	require.Equal(t, nil, val)
}

func TestUpdate(t *testing.T) {
	c := NewExpCache(int64(time.Minute))
	c.Put(1, 1, 1)
	c.Put(1, 2, 1)
	val, ok := c.Get(1, 1)
	require.Equal(t, true, ok)
	require.Equal(t, 2, val)
	require.Equal(t, 1, len(c.data))
	require.Equal(t, 1, c.expL.Len())
}

func TestExpiration(t *testing.T) {
	c := NewExpCache(int64(time.Nanosecond))
	c.Put(1, 1, 1)
	time.Sleep(time.Nanosecond)
	val, ok := c.Get(1, 1)
	require.Equal(t, false, ok)
	require.Equal(t, nil, val)
}

func TestCleanup(t *testing.T) {
	c := NewExpCache(int64(time.Nanosecond))
	c.Put(1, 1, 1)
	c.Put(2, 2, 2)
	time.Sleep(2 * time.Nanosecond)
	require.Equal(t, 2, c.expL.Len())
	require.Equal(t, 2, len(c.data))
	cleaned := c.Cleanup()
	require.Equal(t, 0, c.expL.Len())
	require.Equal(t, 0, len(c.data))
	require.Equal(t, []uint64{1, 2}, cleaned)
}
