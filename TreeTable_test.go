package ux

import (
	"testing"

	"github.com/ddkwork/golibrary/stream"
)

func TestTreeTable_ContextMenuItem(t1 *testing.T) {
	m := safemap.NewOrdered[string, string]()
	m.Set("CopyRow", "CopyRow")
	m.Set("ConvertToContainer", "ConvertToContainer")
	m.Set("ConvertToNonContainer", "ConvertToNonContainer")
	m.Set("New", "New")
	m.Set("NewContainer", "NewContainer")
	m.Set("Delete", "Delete")
	m.Set("Duplicate", "Duplicate")
	m.Set("Edit", "Edit")
	m.Set("OpenAll", "OpenAll")
	m.Set("CloseAll", "CloseAll")
	g := stream.NewGeneratedFile()
	g.EnumTypes("ContextMenuItem", m)
}
