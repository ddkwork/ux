package main

import (
	"demo/test/dirTree/p/gitignore"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ddkwork/golibrary/mylog"
)

// https://github.com/datumbrain/dirtree
func main() {
	var rootDir string

	switch len(os.Args) {
	case 1:
		rootDir = "."
		fmt.Println(".")
	case 2:
		rootDir = os.Args[1]
		mylog.Check2(os.Stat(rootDir))

	default:
		fmt.Println("Usage: dirtree [directory]")
		os.Exit(1)
	}

	printTree(rootDir, "", nil)
}

func shouldIgnore(path string, ignoreMatchers []*gitignore.GitIgnore) bool {
	for _, matcher := range ignoreMatchers {
		if matcher.MatchesPath(path) {
			return true
		}
	}
	return false
}

func loadIgnoreMatchers(path string, parentMatchers []*gitignore.GitIgnore) []*gitignore.GitIgnore {
	ignoreFile := filepath.Join(path, ".gitignore")
	if _, e := os.Stat(ignoreFile); os.IsNotExist(e) {
		return parentMatchers
	}

	ignoreMatcher := mylog.Check2(gitignore.CompileIgnoreFile(ignoreFile))

	return append(parentMatchers, ignoreMatcher)
}

func printTree(root string, prefix string, parentMatchers []*gitignore.GitIgnore) {
	ignoreMatchers := loadIgnoreMatchers(root, parentMatchers)

	files := mylog.Check2(os.ReadDir(root))

	for i, file := range files {
		if file.Name() == ".git" {
			continue
		}

		relPath, _ := filepath.Rel(".", filepath.Join(root, file.Name()))
		if shouldIgnore(relPath, ignoreMatchers) {
			continue
		}

		if i == len(files)-1 {
			fmt.Printf("%s└── %s\n", prefix, file.Name())
		} else {
			fmt.Printf("%s├── %s\n", prefix, file.Name())
		}

		if file.IsDir() {
			newPrefix := prefix
			if i == len(files)-1 {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			printTree(filepath.Join(root, file.Name()), newPrefix, ignoreMatchers)
		}
	}
}
