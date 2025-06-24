package stream

import (
	"embed"
	"path/filepath"

	"github.com/ddkwork/golibrary/std/safemap"

	"github.com/ddkwork/golibrary/std/mylog"
)

func GoReleaser() {
	WriteTruncate(".goreleaser.yaml", NewBuffer(".goreleaser.yaml").String())
	RunCommand("goreleaser release --snapshot")
}

func ReadEmbedFileMap(embedFiles embed.FS, dir string) *safemap.M[string, []byte] {
	fileContents := new(safemap.M[string, []byte])
	fileList := mylog.Check2(embedFiles.ReadDir(dir))
	for _, file := range fileList {
		uncPath := filepath.ToSlash(filepath.Join(dir, file.Name()))
		fileData := mylog.Check2(embedFiles.ReadFile(uncPath))
		fileContents.Set(file.Name(), fileData)
	}
	return fileContents
}

func TrimSlash(name string) string {
	if len(name) > 0 && name[len(name)-1] == '/' {
		return name[:len(name)-1]
	}
	return name
}
