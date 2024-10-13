package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFile(t *testing.T) {
	a := &File{"a", 100, false, []*File{}, "", 0, 0}
	build := NewTestFile("a", 100)
	assert.Equal(t, a, build)
}

func TestBuildFolder(t *testing.T) {
	a := &File{"a", 0, true, []*File{}, "", 0, 0}
	build := NewTestFolder("a")
	assert.Equal(t, a, build)
}

func TestBuildFolderWithFile(t *testing.T) {
	e := &File{"e", 100, false, []*File{}, "", 0, 0}
	d := &File{"d", 100, true, []*File{e}, "", 0, 0}
	build := NewTestFolder("d", NewTestFile("e", 100))
	assert.Equal(t, d, build)
}

func TestBuildComplexFolder(t *testing.T) {
	e := &File{"e", 100, false, []*File{}, "", 0, 0}
	d := &File{"d", 100, true, []*File{e}, "", 0, 0}
	b := &File{"b", 50, false, []*File{}, "", 0, 0}
	c := &File{"c", 100, false, []*File{}, "", 0, 0}
	a := &File{"a", 250, true, []*File{b, c, d}, "", 0, 0}
	build := NewTestFolder("a", NewTestFile("b", 50), NewTestFile("c", 100), NewTestFolder("d", NewTestFile("e", 100)))
	assert.Equal(t, a, build)
}

func TestFindTestFile(t *testing.T) {
	folder := NewTestFolder("a",
		NewTestFolder("b",
			NewTestFile("c", 10),
			NewTestFile("d", 100),
		),
	)
	expected := folder.Files[0].Files[1]
	foundFile := FindTestFile(folder, "d")
	assert.Equal(t, expected, foundFile)
}
