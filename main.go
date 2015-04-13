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
	Event fsnotify.Event

	Filter func(fsnotify.Event) bool

	EventTransmitter struct {
		Events chan Event
		filter Filter
	}

	Watcher struct {
		wMutex    sync.Mutex
		fswatcher *fsnotify.Watcher

		etrans []*EventTransmitter
		Errors chan error
	}
)

var (
	ErrNotADir = errors.New("Provided path is not a directory")
)

func NewWatcher(ets ...*EventTransmitter) (*Watcher, error) {
	fswatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	w := &Watcher{
		fswatcher: fswatcher,
		etrans:    ets,
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
		w.wMutex.Lock()
		err := w.fswatcher.Add(d)
		w.wMutex.Unlock()

		if err != nil {
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
					w.AddDir(fsev.Name)
				}
				for _, e := range w.etrans {
					e.Transmit(fsev)
				}
			case err := <-w.fswatcher.Errors:
				w.Errors <- err
			}
		}
	}()
}
