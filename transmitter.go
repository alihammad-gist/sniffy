package sniffy

import (
	"gopkg.in/fsnotify.v1"
)

// Creates a new filter and makes the channel to be used
// for Event transmission
func Transmitter(filters ...Filter) *EventTransmitter {
	e := make(chan Event)
	return &EventTransmitter{
		filter: FilterChain(filters...),
		Events: e,
	}
}

func (et *EventTransmitter) Transmit(e fsnotify.Event) {
	et.lastEventLock.Lock()
	defer et.lastEventLock.Unlock()

	if et.filter(e) {
		et.Events <- Event(e)
		et.lastEvent = Event(e)
	}
}

func (et *EventTransmitter) LastEvent() Event {
	et.lastEventLock.Lock()
	defer et.lastEventLock.Unlock()

	return et.lastEvent
}
