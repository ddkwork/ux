package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/ddkwork/golibrary/mylog"
)

// https://github.com/aleksandrzaykov88/learngo/tree/refs/heads/master/Coursera/tree

// Leaf is a dir with or without children-cats.
type Leaf struct {
	File   fs.FileInfo
	Leaves []Leaf
}

// Name returning.
func (l Leaf) Name() string {
	if l.File.IsDir() {
		return l.File.Name()
	} else {
		return fmt.Sprintf("%s (%s)", l.File.Name(), l.Size())
	}
}

// Size returning.
func (l Leaf) Size() string {
	if l.File.Size() > 0 {
		return fmt.Sprintf("%db", l.File.Size())
	} else {
		return "empty"
	}
}

// getLeaves gets all leaves from tree.
func getLeaves(path string, printFiles bool) []Leaf {
	files := mylog.Check2(os.ReadDir(path))

	var leaves []Leaf
	for _, file := range files {
		if file.Name() == ".git" {
			continue
		}
		if !printFiles && !file.IsDir() {
			continue
		}
		leaf := Leaf{File: mylog.Check2(file.Info())}
		if file.IsDir() {
			children := getLeaves(path+"\\"+file.Name(), printFiles)
			leaf.Leaves = children
		}
		leaves = append(leaves, leaf)
	}
	return leaves
}

// printLeaves outputs tree in output stream.
func printLeaves(out io.Writer, leaves []Leaf, parentPrefix string) {
	lastIdx := len(leaves) - 1
	prefix := "├───"
	childPrefix := "│\t"
	for i, leaf := range leaves {
		if i == lastIdx {
			prefix = "└───"
			childPrefix = "\t"
		}
		fmt.Fprint(out, parentPrefix, prefix, leaf.Name(), "\n")
		if leaf.File.IsDir() {
			printLeaves(out, leaf.Leaves, parentPrefix+childPrefix)
		}
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	leaves := getLeaves(path, printFiles)
	printLeaves(out, leaves, "")
	return nil
}

func main() {
	out := os.Stdout
	//if !(len(os.Args) == 2 || len(os.Args) == 3) {
	//	panic("usage go run main.go . [-f]")
	//}
	//path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	printFiles = true
	mylog.Check(dirTree(out, ".", printFiles))
}
