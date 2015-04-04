package sniffyTest

import (
	"path/filepath"

	"os"
)

func getDir() {
	dir := filepath.Join(
		os.TempDir(),
		"sniffy_test",
	)
	os.Remove(dir)
	os.Create(dir) // incorrect! probably is for files only MkdirAll <-
	return dir
}
