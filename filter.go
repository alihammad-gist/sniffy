package sniffy

import "gopkg.in/fsnotify.v1"

// This factory will concatinate multiple filters into
// one
func FilterChain(fs ...Filter) Filter {
	return func(e fsnotify.Event) bool {
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

func AutoWatch(w Watcher) Filter {
	return func(e fsnotify.Event) bool {
		if e.Op == fsnotify.Create {
			if isDir(e.Name) {
				w.AddDir(e.Name)
			}
		}
		return true
	}
}
