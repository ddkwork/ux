package ignore

import (
	"bufio"

	"os"
	"os/user"
	"path/filepath"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux/filetree/files"
)

func ReadIgnoreFile() []string {
	usr := mylog.Check2(user.Current())

	ignoreFileName := filepath.Join(usr.HomeDir, ".goduignore")
	if _ := mylog.Check2(os.Stat(ignoreFileName)); os.IsNotExist(err) {
		return []string{}
	}
	ignoreFile := mylog.Check2(os.Open(ignoreFileName))

	defer ignoreFile.Close()
	scanner := bufio.NewScanner(ignoreFile)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func IgnoreBasedOnIgnoreFile(ignoreFile []string) files.ShouldIgnoreFolder {
	ignoredFolders := map[string]struct{}{}
	for _, line := range ignoreFile {
		ignoredFolders[line] = struct{}{}
	}
	return func(absolutePath string) bool {
		_, name := filepath.Split(absolutePath)
		_, ignored := ignoredFolders[name]
		return ignored
	}
}
