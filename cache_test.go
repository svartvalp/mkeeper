package mkeeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEviction(t *testing.T) {
	c := NewCache(WithTTL(time.Second))
	c.Put(1, 1)
	time.Sleep(2 * time.Second)
	val, ok := c.GetIfPresent(1)
	require.Equal(t, false, ok)
	require.Equal(t, nil, val)
}

func TestMaxCapacity(t *testing.T) {
	c := NewCache(WithMaxCapacity(1))
	c.Put(1, 1)
	c.Put(2, 2)
	val, ok := c.GetIfPresent(2)
	require.Equal(t, false, ok)
	require.Equal(t, nil, val)
}

type testStruct struct {
	A string
	B int
	C byte
}

func TestStructKey(t *testing.T) {
	str := testStruct{
		A: "1",
		B: 2,
		C: 1,
	}
	str2 := testStruct{
		A: "1",
		B: 2,
		C: 1,
	}
	c := NewCache()
	c.Put(str, 1)
	c.Put(str, 2)
	val, ok := c.GetIfPresent(str2)
	require.Equal(t, true, ok)
	require.Equal(t, 2, val)
}

func TestBufferOverflow(t *testing.T) {
	c := NewCache(WithTTL(time.Millisecond), WithCleanBuf(1), WithCleanTick(time.Second/2))
	c.Put(1, 1)
	c.Put(2, 2)
	time.Sleep(time.Second)
	val, ok := c.GetIfPresent(1)
	require.Equal(t, false, ok)
	require.Equal(t, nil, val)
}

func TestAdmissionPolicy(t *testing.T) {
	c := NewCache(WithTTL(time.Minute), WithMaxCapacity(2))
	c.Put(1, 1)
	c.Put(2, 2)
	c.GetIfPresent(2)
	c.GetIfPresent(1)
	time.Sleep(time.Second)
	c.Put(3, 3)
	val, ok := c.GetIfPresent(3)
	require.Equal(t, false, ok)
	require.Equal(t, nil, val)
	val, ok = c.GetIfPresent(2)
	require.Equal(t, true, ok)
	require.Equal(t, 2, val)
	val, ok = c.GetIfPresent(1)
	require.Equal(t, true, ok)
	require.Equal(t, 1, val)
}

func TestAdmissionPolicy_Victim(t *testing.T) {
	c := NewCache(WithTTL(time.Minute), WithMaxCapacity(2))
	c.Put(1, 1)
	c.Put(2, 2)
	c.GetIfPresent(3)
	c.GetIfPresent(3)
	c.GetIfPresent(1)
	time.Sleep(time.Second)
	c.Put(3, 3)
	val, ok := c.GetIfPresent(3)
	require.Equal(t, true, ok)
	require.Equal(t, 3, val)
}
