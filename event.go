package mkeeper

import (
	"time"
)

type EventListener struct {
	processors []EventProcessor

	ch     chan Event
	buf    []Event
	maxBuf int
	tick   *time.Ticker
}

type EventProcessor interface {
	ProcessEvents(e []Event)
}

func NewEventListener(procs []EventProcessor, maxBuf int, tickDur time.Duration) *EventListener {
	return &EventListener{
		processors: procs,
		ch:         make(chan Event, 2048),
		buf:        make([]Event, 0, maxBuf),
		maxBuf:     maxBuf,
		tick:       time.NewTicker(tickDur),
	}
}

type Event struct {
	Type EventType
	Hash uint64
}

type EventType int

const (
	Put    EventType = 1
	Delete EventType = 2
	Get    EventType = 3
	Miss   EventType = 4
)

func (e *EventListener) Send(event Event) {
	e.ch <- event
}

func (e *EventListener) Start() {
	go func() {
		for {
			select {
			case val := <-e.ch:
				e.buf = append(e.buf, val)
				if len(e.buf) >= e.maxBuf {
					for _, p := range e.processors {
						p.ProcessEvents(e.buf)
					}
					e.buf = make([]Event, 0, e.maxBuf)
				}
			case <-e.tick.C:
				for _, p := range e.processors {
					p.ProcessEvents(e.buf)
				}
				e.buf = make([]Event, 0, e.maxBuf)
			}
		}
	}()
}
