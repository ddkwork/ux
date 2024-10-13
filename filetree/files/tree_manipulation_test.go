package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortFolder(t *testing.T) {
	folder := NewTestFolder("b",
		NewTestFile("c", 100),
		NewTestFolder("d",
			NewTestFile("e", 50),
			NewTestFile("f", 30),
			NewTestFolder("g",
				NewTestFile("i", 30),
				NewTestFile("j", 50),
			),
		),
	)
	expected := NewTestFolder("b",
		NewTestFolder("d",
			NewTestFolder("g",
				NewTestFile("j", 50),
				NewTestFile("i", 30),
			),
			NewTestFile("e", 50),
			NewTestFile("f", 30),
		),
		NewTestFile("c", 100),
	)

	SortDesc(folder)
	assert.Equal(t, expected, folder)
}

func TestPruneFolder(t *testing.T) {
	folder := &File{"b", 260, true, []*File{
		{"c", 100, false, []*File{}, "", 1, 0},
		{"d", 160, true, []*File{
			{"e", 50, false, []*File{}, "", 2, 0},
			{"f", 30, false, []*File{}, "", 2, 0},
			{"g", 80, true, []*File{
				{"i", 50, false, []*File{}, "", 3, 0},
				{"j", 30, false, []*File{}, "", 3, 0},
			}, "", 2, 0},
		}, "", 1, 0},
	}, "", 0, 0}
	expected := &File{"b", 260, true, []*File{
		{"c", 100, false, []*File{}, "", 1, 0},
		{"d", 160, true, []*File{
			{"g", 80, true, []*File{}, "", 2, 0},
		}, "", 1, 0},
	}, "", 0, 0}
	PruneSmallFiles(folder, 60)
	assert.Equal(t, expected, folder)
}
