package mkeeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testEvProc struct {
	f func(ev []Event)
}

func (t *testEvProc) ProcessEvents(e []Event) {
	t.f(e)
}

func TestEventListener(t *testing.T) {
	ev := NewEventListener([]EventProcessor{&testEvProc{
		f: func(ev []Event) {
			if len(ev) > 0 {
				require.Equal(t, uint64(1), ev[0].Hash)
				require.Equal(t, uint64(2), ev[1].Hash)
				require.Equal(t, uint64(3), ev[2].Hash)
			}
		},
	}}, 3, time.Second)

	ev.Start()
	ev.Send(Event{
		Type: Put,
		Hash: 1,
	})
	ev.Send(Event{
		Type: Delete,
		Hash: 2,
	})
	ev.Send(Event{
		Type: Put,
		Hash: 3,
	})
	time.Sleep(time.Second / 6)
}
