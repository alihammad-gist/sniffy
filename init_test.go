package sniffy_test

import (
	"log"
	"path/filepath"

	"os"
)

var Paths = map[string][]string{
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

func getDir() {
	tmp := os.TempDir()
	for d, fs := range Dirs {
		path := filepath.Join(tmp, d)
		os.Remove(path)
		if err := os.MkdirAll(path, os.ModeDir); err != nil {
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
