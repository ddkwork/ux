package ux

import (
	"github.com/ddkwork/golibrary/stream"
	"github.com/goradd/maps"
	"testing"
)

func TestTreeTable_ContextMenuItem(t1 *testing.T) {
	m := new(maps.SafeSliceMap[string, string])
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
