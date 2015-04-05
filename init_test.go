package sniffy_test

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/alihammad-gist/sniffy"
	"gopkg.in/fsnotify.v1"

	"os"
)

var Paths = map[string][]string{
	"sniffy_test": {
		"root.tree",
	},
	"sniffy_test/l1d1": {
		"1.txt",
	},
	"sniffy_test/l1d1/l2d1": {
		"index.php",
		"default.css",
	},
	"sniffy_test/l1d1/l2d2": {
		"main.cpp",
	},
	"sniffy_test/l1d2": {
		"fileserver.py",
	},
	"sniffy_test/l1d2/l2d1": {
		"hello",
	},
	"sniffy_test/l1d2/l2d2": {
		"main.less",
		"vars.less",
	},
	"sniffy_test/l1d3": {
		"init.cfx",
		"global.log",
	},
	"sniffy_test/l1d3/l2d1": {
		"main.sass",
		"colors.scss",
	},
	"sniffy_test/l1d3/l2d2": {
		"pre.xstream",
	},
}

func getDir() string {
	tmp := os.TempDir()
	os.RemoveAll(filepath.Join(tmp, "sniffy_test"))
	for d, fs := range Paths {
		path := filepath.Join(tmp, d)
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatal(err)
		}
		for _, f := range fs {
			if _, err := os.Create(
				filepath.Join(path, f),
			); err != nil {
				log.Fatal(err)
			}
		}
	}
	return filepath.Join(tmp, "sniffy_test")
}

func triggerOperation(path string, op fsnotify.Op) {
	switch op {
	case sniffy.Create:
		if _, err := os.Create(path); err != nil {
			log.Fatal("Triggering C-operation error:", err)
		}
	case sniffy.Remove:
		if err := os.RemoveAll(path); err != nil {
			log.Fatal("Triggering R-operation error:", err)
		}
	case sniffy.Rename:
		if err := os.Rename(
			path,
			filepath.Join(filepath.Dir(path), "renamed.txt"),
		); err != nil {
			log.Fatal("Triggering Rn-operation error:", err)
		}
	case sniffy.Write:
		if err := ioutil.WriteFile(
			path,
			[]byte("hello@!#"),
			0755,
		); err != nil {
			log.Fatal("Triggering w-operation error:", err)
		}
	}
}
