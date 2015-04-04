package sniffy_test

import (
	"log"

	"testing"
)

func TestWatcherWithNoFilters(t *testing.T) {
	dir := getDir()
	log.Println("Dirs created: ", dir)
}
