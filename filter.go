package sniffy

import (
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/fsnotify.v1"
)

// This factory will concatinate multiple filters into
// one
func FilterChain(fs ...Filter) Filter {
	m := &sync.Mutex{}
	return func(e fsnotify.Event) bool {
		m.Lock()
		defer m.Unlock()
		for _, f := range fs {
			if !f(e) {
				return false
			}
		}
		return true
	}
}

// If Event was triggered by provided operation this factory
// will return Filter that will pass
func OpFilter(ops ...fsnotify.Op) Filter {
	return func(e fsnotify.Event) bool {
		for _, op := range ops {
			if op == e.Op {
				return true
			}
		}
		return false
	}
}

// Returns true only if event occured on files with
// provided extensions
func ExtFilter(exts ...string) Filter {
	return func(e fsnotify.Event) bool {
		for _, ext := range exts {
			if filepath.Ext(e.Name) == ext {
				return true
			}
		}
		return false
	}
}

// Returns true only if event occured on a child
// of provided paths
// paths must be absolute
func ChildFilter(paths ...string) Filter {
	return func(e fsnotify.Event) bool {
		for _, p := range paths {
			if strings.HasPrefix(e.Name+"/", p) {
				return true
			}
		}
		return false
	}
}

// Returns false if Event path is one of the provided paths
// paths must be absolute
func ExcludePathFilter(paths ...string) Filter {
	return func(e fsnotify.Event) bool {
		for _, p := range paths {
			if p == e.Name {
				return false
			}
		}
		return true
	}
}

// Returns false if last event occured
// within the specified duration
// It can be used to filter out simultaneous events
func TooSoonFilter(d time.Duration) Filter {
	var (
		lastTime time.Time
	)
	return func(_ fsnotify.Event) bool {
		now := time.Now()
		if !lastTime.IsZero() {
			if now.Sub(lastTime) <= d {
				return false
			}
		}
		lastTime = now
		return true
	}
}

// Returns false if event occured on file with provided
// shell filename patterns, pattern will be matched against
// filename not the absolute path
func IgnoreFnPatternFilter(pats ...string) Filter {
	return func(e fsnotify.Event) bool {
		for _, p := range pats {
			if ok, _ := filepath.Match(p, filepath.Base(e.Name)); ok {
				return false
			}
		}
		return true
	}
}
