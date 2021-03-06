package sniffy_test

import (
	"testing"
	"time"

	"github.com/alihammad-gist/sniffy"
	"gopkg.in/fsnotify.v1"
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
		{sniffy.Event{"/home/ali/ha/121313", sniffy.Chmod}, false},
		{sniffy.Event{"/home/ali/.index.php~", sniffy.Chmod}, false},
	}

	for _, ev := range evs {
		if extf(fsnotify.Event(ev.e)) != ev.x {
			t.Logf("Expected %t Event %v", ev.x, ev.e)
			t.Fail()
		}
	}
}

func TestExcludeChild(t *testing.T) {
	childf := sniffy.ChildFilter("/name/app", "/usr/bin")
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
		if childf(fsnotify.Event(ev.e)) != ev.x {
			t.Logf("Expected %t Event %v", ev.x, ev.e)
			t.Fail()
		}
	}
}

func TestExcludePathFilter(t *testing.T) {
	pathf := sniffy.ExcludePathFilter("/home/ali/bin", "/usr/bin/ali")
	evs := []struct {
		e sniffy.Event
		x bool
	}{
		{sniffy.Event{"/home/ali.php", sniffy.Chmod}, true},
		{sniffy.Event{"/name/app/vars.sass", sniffy.Chmod}, true},
		{sniffy.Event{"/home/hello/main.css", sniffy.Chmod}, true},
		{sniffy.Event{"/home/ali.php/main.txt", sniffy.Chmod}, true},
		{sniffy.Event{"/home/ali/bin", sniffy.Chmod}, false},
		{sniffy.Event{"/usr/bin/ali", sniffy.Chmod}, false},
	}
	for _, ev := range evs {
		if pathf(fsnotify.Event(ev.e)) != ev.x {
			t.Logf("Expected %t Event %v", ev.x, ev.e)
			t.Fail()
		}
	}
}

func TestTooSoonFilter(t *testing.T) {
	soonf := sniffy.TooSoonFilter(time.Millisecond * 500)
	evs := []struct {
		e sniffy.Event
		d time.Duration
		x bool
	}{
		{sniffy.Event{"/path/1", sniffy.Chmod}, time.Millisecond, true},
		{sniffy.Event{"/path/1", sniffy.Chmod}, time.Millisecond * 501, true},
		{sniffy.Event{"/path/1", sniffy.Chmod}, time.Millisecond * 500, true},
		{sniffy.Event{"/path/2", sniffy.Chmod}, time.Millisecond, false},
		{sniffy.Event{"/path/1", sniffy.Chmod}, time.Millisecond, false},
		{sniffy.Event{"/path/2", sniffy.Chmod}, time.Millisecond, false},
		{sniffy.Event{"/path/2", sniffy.Chmod}, time.Millisecond, false},
		{sniffy.Event{"/path/1", sniffy.Chmod}, time.Millisecond, false},
		{sniffy.Event{"/path/1", sniffy.Chmod}, time.Millisecond, false},
		{sniffy.Event{"/path/1", sniffy.Create}, time.Millisecond * 400, false},
		{sniffy.Event{"/path/1", sniffy.Create}, time.Millisecond * 500, true},
	}

	for _, ev := range evs {
		<-time.After(ev.d)
		if soonf(fsnotify.Event(ev.e)) != ev.x {
			t.Logf("Expected %t Event %v Duration %v", ev.x, ev.e, ev.d)
			t.Fail()
		}
	}
}

func TestIgnoreFnPatternFilter(t *testing.T) {
	var (
		evs = []struct {
			e sniffy.Event
			x bool
		}{
			{sniffy.Event{"/path1/2.txt", sniffy.Chmod}, true},
			{sniffy.Event{"/path324/.3.txt", sniffy.Chmod}, false},
			{sniffy.Event{"/a/b/c/d.php", sniffy.Chmod}, false},
			{sniffy.Event{"/a/b/c.hello", sniffy.Chmod}, true},
			{sniffy.Event{"x/y/.z.cpp", sniffy.Chmod}, false},
		}

		pat = []string{".??[^.]*", "*.php"}

		fil = sniffy.IgnoreFnPatternFilter(pat...)
	)

	for _, ev := range evs {
		if fil(fsnotify.Event(ev.e)) != ev.x {
			t.Errorf("IgnoreFnPatterFilter: Filename  %s, Expected  %t", ev.e.Name, ev.x)
			t.Fail()
		}
	}
}
