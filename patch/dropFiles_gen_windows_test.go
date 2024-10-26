package app

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ddkwork/golibrary/stream"
)

// SetTitleBarIsDark
// genDropFilesEventArg
func TestGenDropFiles(t *testing.T) {
	stream.CopyFile("fix/os_windows.go", "../vendor/gioui.org/app/os_windows.go")
	stream.CopyFile("fix/dropFiles_gen_windows.go", "../vendor/gioui.org/app/dropFiles_gen_windows.go")
	for _, s := range getOsFiles() {
		if filepath.Base(s) == "os_windows.go" {
			continue
		}
		println(s)
		b := stream.NewBuffer(s)
		if b.Contains(callbackDefault) {
			continue
		}
		stream.WriteAppend(s, callbackDefault)
	}
	stream.NewBuffer("../vendor/gioui.org/x/styledtext/styledtext.go").ReplaceAll("MaxWidth:   maxWidth,", "MaxWidth:   maxWidth*2,").ReWriteSelf()

	stream.CopyFile("apk/androidbuild.go", "../vendor/gioui.org/cmd/gogio/androidbuild.go")
	os.Chdir("../vendor/gioui.org/cmd/gogio")
	stream.RunCommand("go install .")
}

var callbackDefault = `
var dragHandler = func(files []string) {}

func FileDropCallback(fn func(files []string)) {
	dragHandler = fn
}
`

func getOsFiles() []string {
	var paths []string
	filepath.Walk("../vendor/gioui.org/app", func(path string, info fs.FileInfo, err error) error {
		if strings.HasPrefix(filepath.Base(path), "os_") {
			if filepath.Ext(path) == ".go" {
				paths = append(paths, path)
			}
		}
		return err
	})
	return paths
}
