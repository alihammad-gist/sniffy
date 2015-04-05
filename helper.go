package sniffy

import (
	"os"
	"path/filepath"
)

func isDir(path string) bool {
	fInfo, err := os.Stat(path)
	if err == nil && fInfo.IsDir() {
		return true
	}
	return false
}

func dirTree(root string) chan string {
	dir := make(chan string)
	go func() {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				dir <- path
			}
			return nil
		})
		close(dir)
	}()
	return dir
}
