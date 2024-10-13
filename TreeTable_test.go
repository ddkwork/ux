package ux

import (
	"github.com/ddkwork/golibrary/stream"
	"testing"
)

func TestTreeTable_ContextMenuItem(t1 *testing.T) {
	stream.NewGeneratedFile().Enum("ContextMenuItem",
		[]string{
			"CopyRow",
			"ConvertToContainer",
			"ConvertToNonContainer",
			"New",
			"NewContainer",
			"Delete",
			"Duplicate",
			"Edit",
			"OpenAll",
			"CloseAll",
		},
		nil)
}
