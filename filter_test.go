package sniffy_test

import (
	"github.com/alihammad-gist/sniffy"
	"gopkg.in/fsnotify.v1"
	"testing"
)

func TestOpFilter(t *testing.T) {
	opf := sniffy.OpFilter(sniffy.Remove, sniffy.Chmod)
	evs := []struct {
		e sniffy.Event
		x bool
	}{
		{sniffy.Event{"...", sniffy.Chmod}, true},
		{sniffy.Event{"...", sniffy.Remove}, true},
		{sniffy.Event{"...", sniffy.Create}, false},
		{sniffy.Event{"...", sniffy.Rename}, false},
	}

	for _, ev := range evs {
		if opf(fsnotify.Event(ev.e)) != ev.x {
			t.Logf("Expected %t Event %v", ev.x, ev.e)
			t.Fail()
		}
	}
}

func TestExtFilter(t *testing.T) {
	extf := sniffy.ExtFilter(".php", ".css", ".sass")
	evs := []struct {
		e sniffy.Event
		x bool
	}{
		{sniffy.Event{"/home/ali.php", sniffy.Chmod}, true},
		{sniffy.Event{"app/vars.sass", sniffy.Chmod}, true},
		{sniffy.Event{"/home/hello/main.css", sniffy.Chmod}, true},
		{sniffy.Event{"/home/ali.php/main.txt", sniffy.Chmod}, false},
		{sniffy.Event{"/home/ali", sniffy.Chmod}, false},
	}

	for _, ev := range evs {
		if extf(fsnotify.Event(ev.e)) != ev.x {
			t.Logf("Expected %t Event %v", ev.x, ev.e)
			t.Fail()
		}
	}
}

func TestPathFilter(t *testing.T) {
	pathf := sniffy.PathFilter("/name/app", "/usr/bin")
	evs := []struct {
		e sniffy.Event
		x bool
	}{
		{sniffy.Event{"/home/ali.php", sniffy.Chmod}, false},
		{sniffy.Event{"/name/app/vars.sass", sniffy.Chmod}, true},
		{sniffy.Event{"/home/hello/main.css", sniffy.Chmod}, false},
		{sniffy.Event{"/home/ali.php/main.txt", sniffy.Chmod}, false},
		{sniffy.Event{"/home/ali/bin", sniffy.Chmod}, false},
		{sniffy.Event{"/usr/bin/ali", sniffy.Chmod}, true},
	}
	for _, ev := range evs {
		if pathf(fsnotify.Event(ev.e)) != ev.x {
			t.Logf("Expected %t Event %v", ev.x, ev.e)
			t.Fail()
		}
	}
}
