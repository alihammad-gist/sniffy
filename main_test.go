package sniffy_test

import (
	"os"
	"path/filepath"

	"gopkg.in/fsnotify.v1"

	"github.com/alihammad-gist/sniffy"

	"testing"
	"time"
)

func TestRecursiveWatch(t *testing.T) {
	dir := getDir()
	w, err := sniffy.NewWatcher()
	defer w.Close()
	go func() {

	}()
	if err != nil {
		t.Log("Error encountered:", err)
		t.Fail()
	}
	w.AddDir(dir)
	bye := filepath.Join(dir, "l1d1/l2d1/bye.txt")

	if _, err = os.Create(bye); err != nil {
		t.Log("Error encountered:", err)
		t.Fail()
	}
	select {
	case e := <-w.Events:
		if e.Name != bye || e.Op != sniffy.Create {
			t.Log("Expecting ", bye, "; Found ", e.Name)
			t.Fail()
		}
	case er := <-w.Errors:
		t.Log("Event Error:", er)
		t.Fail()
	case <-time.After(time.Second * 2):
		t.Log("Event taking too long")
		t.Fail()
	}

	if err = os.Remove(bye); err != nil {
		t.Log("Error Removing:", err)
		t.Fail()
	}
	select {
	case e := <-w.Events:
		if e.Name != bye || e.Op != sniffy.Remove {
			t.Log("Expecting ", bye, "; Found ", e.Name)
			t.Fail()
		}
	case er := <-w.Errors:
		t.Log("Event Error:", er)
		t.Fail()
	case <-time.After(time.Second * 2):
		t.Log("Event taking too long")
		t.Fail()
	}
}

func TestRecurWatch(t *testing.T) {
	dir := getDir()
	ops := map[string]fsnotify.Op{
		filepath.Join(dir, "root.tree"):             sniffy.Remove,
		filepath.Join(dir, "root.tree"):             sniffy.Create,
		filepath.Join(dir, "l1d1/l2d1/hulu.conf"):   sniffy.Create,
		filepath.Join(dir, "l1d3/global.log"):       sniffy.Write,
		filepath.Join(dir, "l1d3/l2d2/pre.xstream"): sniffy.Rename,
		filepath.Join(dir, "l1d2/l2d2"):             sniffy.Remove,
		filepath.Join(dir, "l1d2/l2d2"):             sniffy.Create,
	}
	w, err := sniffy.NewWatcher()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer w.Close()
	for path, op := range ops {
		triggerOperation(path, op)
		select {
		case e := <-w.Events:
			if e.Name != path || e.Op != op {
				t.Log("Expecting:", path, "Found:", e.Name)
				t.Fail()
			}
		case er := <-w.Errors:
			t.Log("Event Error:", er)
			t.Fail()
		case <-time.After(time.Second * 2):
			t.Log("Event taking too long")
			t.Fail()
		}
	}
}
