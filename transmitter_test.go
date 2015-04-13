package sniffy_test

import (
	"github.com/alihammad-gist/sniffy"
	"gopkg.in/fsnotify.v1"
	"testing"
	"time"
)

func TestTransmitter(t *testing.T) {
	done := make(chan bool)
	trans := sniffy.Transmitter(func(e fsnotify.Event) bool {
		return false
	})
	go func() {
		select {
		case e := <-trans.Events:
			t.Logf("Unexpected event capture %v", e)
			t.Fail()
			done <- true
		case <-time.After(time.Second):
			done <- true
		}
	}()
	trans.Transmit(fsnotify.Event{})
	<-done

	trans = sniffy.Transmitter(func(e fsnotify.Event) bool {
		return true
	})
	go func() {
		select {
		case <-trans.Events:
			done <- true
		case <-time.After(time.Second):
			t.Logf("Couldn't capture event")
			t.Fail()
			done <- true
		}
	}()
	trans.Transmit(fsnotify.Event{})
	<-done
}
