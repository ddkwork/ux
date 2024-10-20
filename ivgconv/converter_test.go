package ivgconv

import (
	"bytes"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux/ivgconv/testdata"
)

func TestFromFile(t *testing.T) {
	// Encode the SVG file as IconVG.
	ivgData := mylog.Check2(FromFile("testdata/close.svg"))

	// Check that the IconVG data matches the expected output.
	if !bytes.Equal(ivgData, testdata.Close) {
		t.Fatalf("ivgData != Close")
	}

	// Encode the SVG file as IconVG.
	ivgData = mylog.Check2(FromFile("testdata/StarHalf.svg"))

	// Check that the IconVG data matches the expected output.
	if !bytes.Equal(ivgData, testdata.StarHalf) {
		t.Fatalf("ivgData != StarHalf")
	}
}
