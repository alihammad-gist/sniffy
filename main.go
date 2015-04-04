package sniffy

import (
	"errors"
	"log"

	"gopkg.in/fsnotify.v1"
)

const (
	Create = fsnotify.Create
	Write  = fsnotify.Write
	Remove = fsnotify.Remove
	Rename = fsnotify.Rename
	Chmod  = fsnotify.Chmod
)

type (
	Filter func(fsnotify.Event) bool

	Event fsnotify.Event

	Watcher struct {
		fswatcher *fsnotify.Watcher
		filter    Filter
		Event     chan Event
		Error     chan error
	}
)

var (
	ErrNotADir = errors.New("Provided path is not a directory")
)

func NewWatcher(filters ...Filter) (*Watcher, error) {
	fswatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	w := &Watcher{
		fswatcher: fswatcher,
		filter:    FilterChain(filters...),
	}
	w.watch()
	return w
}

func (w *Watcher) AddDir(path string) error {
	if !isDir(path) {
		return ErrNotADir
	}
	for d := range w.dirTree(path) {
		w.fswatcher.Add(d)
	}
	return nil
}

func (w *Watcher) Close() error {
	return w.fswatcher.Close()
}

func (w *Watcher) watch() {
	go func() {
		for {
			select {
			case ev := <-w.fswatcher.Events:
				if w.filter(ev) {
					w.Event <- Event(ev)
				}
			case err := <-w.fswatcher.Errors:
				log.Println("Error: ", err)
				w.Error <- err
			}
		}
	}()
}
