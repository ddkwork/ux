package ux

import (
	"testing"

	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/golibrary/stream"
)

func TestTreeTable_ContextMenuItem(t1 *testing.T) {
	m := safemap.NewOrdered[string, string](func(yield func(string, string) bool) {
		yield("CopyRow", "CopyRow")
		yield("ConvertToContainer", "ConvertToContainer")
		yield("ConvertToNonContainer", "ConvertToNonContainer")
		yield("New", "New")
		yield("NewContainer", "NewContainer")
		yield("Delete", "Delete")
		yield("Duplicate", "Duplicate")
		yield("Edit", "Edit")
		yield("OpenAll", "OpenAll")
		yield("CloseAll", "CloseAll")
		yield("SaveData", "SaveData")
	})
	stream.NewGeneratedFile().EnumTypes("ContextMenuItem", m)
}
