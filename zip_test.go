package simplezip_test

import (
	"testing"

	"github.com/ricochhet/simplezip"
)

func TestZip(t *testing.T) { //nolint:paralleltest // dependant
	if err := simplezip.Zip("./", "./.test/simplezip-src.zip"); err != nil {
		t.Fatal(err)
	}
}

func TestUnzip(t *testing.T) { //nolint:paralleltest // dependant
	if err := simplezip.Unzip("./.test/simplezip-src.zip", "./.test/src"); err != nil {
		t.Fatal(err)
	}
}
