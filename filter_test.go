package sniffy_test

import (
	"github.com/alihammad-gist/sniffy"
	"testing"
)

func TestOpFilter(t *testing.T) {
	opf := sniffy.OpFilter(sniffy.Remove, sniffy.Chmod)

}
