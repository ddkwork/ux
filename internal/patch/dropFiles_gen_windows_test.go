package app

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ddkwork/golibrary/std/stream"
)

func TestGenDropFiles(t *testing.T) {
	for _, path := range getOsFiles() {
		if filepath.Base(path) == "os_windows.go" {
			continue
		}
		println(path)
		b := stream.NewBuffer(path)
		if b.Contains(callbackDefault) {
			continue
		}
		stream.WriteAppend(path, callbackDefault)
	}
}

var callbackDefault = `
var dragHandler = func(files []string) {}

func FileDropCallback(fn func(files []string)) {
	dragHandler = fn
}
`

func getOsFiles() []string {
	var paths []string
	filepath.Walk("../../vendor/gioui.org/app", func(path string, info fs.FileInfo, err error) error {
		if strings.HasPrefix(filepath.Base(path), "os_") {
			if filepath.Ext(path) == ".go" {
				paths = append(paths, path)
			}
		}
		return err
	})
	return paths
}
