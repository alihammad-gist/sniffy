package sniffy

import (
	"errors"
	"sync"

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
		wMutex    sync.Mutex
		fswatcher *fsnotify.Watcher

		filter Filter
		Events chan Event
		Errors chan error
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
		Events:    make(chan Event),
		Errors:    make(chan error),
	}
	w.watch()
	return w, nil
}

func (w *Watcher) AddDir(path string) error {
	if !isDir(path) {
		return ErrNotADir
	}
	for d := range dirTree(path) {
		if err := w.fswatcher.Add(d); err != nil {
			return err
		}
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
			case fsev := <-w.fswatcher.Events:
				if isDir(fsev.Name) {
					w.wMutex.Lock()
					w.AddDir(fsev.Name)
					w.wMutex.Unlock()
				}
				if w.filter(fsev) {
					w.Events <- Event(fsev)
				}
			case err := <-w.fswatcher.Errors:
				w.Errors <- err
			}
		}
	}()
}
