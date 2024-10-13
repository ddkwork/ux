package ignore

import (
	"bufio"
	"github.com/ddkwork/ux/filetree/files"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func ReadIgnoreFile() []string {
	usr, err := user.Current()
	if err != nil {
		log.Println("Wasn't able to retrieve current user at runtime")
		return []string{}
	}
	ignoreFileName := filepath.Join(usr.HomeDir, ".goduignore")
	if _, err := os.Stat(ignoreFileName); os.IsNotExist(err) {
		return []string{}
	}
	ignoreFile, err := os.Open(ignoreFileName)
	if err != nil {
		log.Printf("Failed to read ingorefile because %s\n", err.Error())
		return []string{}
	}
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
