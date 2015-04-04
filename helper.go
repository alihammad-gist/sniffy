package sniffy

import "os"

func isDir(path string) bool {
	fInfo, err := os.Stat(path)
	if err == nil && fInfo.IsDir() {
		return true
	}
	return false
}
