package sniffy_test

import (
	"path/filepath"
	"time"

	"github.com/alihammad-gist/sniffy"

	"testing"
)

func TestRecurWatch(t *testing.T) {
	dir := getDir()
	ops := []Op{
		Op{filepath.Join(dir, "root.tree"), sniffy.Remove},
		Op{filepath.Join(dir, "l1d3/global.log"), sniffy.Write},
		Op{filepath.Join(dir, "l1d3/l2d2/pre.xstream"), sniffy.Rename},
		Op{filepath.Join(dir, "l1d2/l2d2/vars.less"), sniffy.Remove},
		Op{filepath.Join(dir, "l1d2/l2d2/vars.less"), sniffy.Create},
		Op{filepath.Join(dir, "l1d1/l2d2/l3d1"), sniffy.Create},
		Op{filepath.Join(dir, "l1d1/l2d2/l3d1/tree.log"), sniffy.Create},
		Op{filepath.Join(dir, "l1d1/l2d2/l3d1/tree.log"), sniffy.Remove},
	}
	for _, xOp := range ops {
		w, err := sniffy.NewWatcher()
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		w.AddDir(dir)
		triggerOperation(xOp.path, xOp.op)
		select {
		case e := <-w.Events:
			if e.Name != xOp.path || e.Op != xOp.op {
				t.Log("Expected", xOp.path, xOp.op, "Actual", e.Name, e.Op)
				t.Fail()
			}
		case err := <-w.Errors:
			t.Log("Watch Error", err, xOp.path, xOp.op)
			t.Fail()
		case <-time.After(WaitDuration):
			t.Log("Timedout", xOp.path, xOp.op)
			t.Fail()
		}
		w.Close()
	}
}
