package ux

import (
	"github.com/ddkwork/golibrary/stream"
	"testing"
)

func TestTreeTable_ContextMenuItem(t1 *testing.T) {
	stream.NewGeneratedFile().Types("ContextMenuItem",
		[]string{
			"CopyRow",
			"ConvertToContainer",
			"ConvertToNonContainer",
			"NewOrderedMap",
			"NewContainer",
			"Delete",
			"Duplicate",
			"Edit",
			"OpenAll",
			"CloseAll",
		},
		nil)
}
