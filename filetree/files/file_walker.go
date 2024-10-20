package files

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"

	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
)

// File structure representing files and folders with their accumulated sizes
type File struct {
	Name        string  // Name of the file
	Size        int64   // Size of the file or directory
	IsDir       bool    // To indicate if the file is a folder or not
	Files       []*File // Files that contain in case IsDir == true.
	FullPath    string  // Path of the file
	Depth       int     // Indicates in which level the file is compared with the root level
	NumChildren int64   // Num of files that the directory contains
}

// This allows to reduce RAM usage. We only have widget.Bools for shown files instead of the whole filesystem
type FileShow struct {
	File         *File       // Points to a file
	IsSelected   widget.Bool // Indicate if the file has been selected
	ActionButton widget.Bool // Indicate if the folder has to be opened/closed
}

// UpdateSize goes through subfiles and subfolders and accumulates their size
func (f *File) UpdateSize(level int) {
	if !f.IsDir {
		return
	}
	var size int64
	var numchildren int64
	for _, child := range f.Files {
		child.UpdateSize(level + 1)
		size += child.Size
		if child.IsDir {
			numchildren += child.NumChildren
		} else {
			numchildren++
		}
	}
	f.Size = size
	f.Depth = level
	f.NumChildren = numchildren

	// Sort files
	sort.Slice(f.Files, func(i, j int) bool {
		return f.Files[i].Size > f.Files[j].Size
	})
}

// ReadDir function can return list of files for given folder path
type ReadDir func(dirname string) ([]os.FileInfo, error)

// ShouldIgnoreFolder function decides whether a folder should be ignored
type ShouldIgnoreFolder func(absolutePath string) bool

func ignoringReadDir(shouldIgnore ShouldIgnoreFolder, originalReadDir ReadDir) ReadDir {
	return func(path string) ([]os.FileInfo, error) {
		if shouldIgnore(path) {
			return []os.FileInfo{}, nil
		}
		return originalReadDir(path)
	}
}

// WalkFolder will go through a given folder and subfolders and produces file structure
// with aggregated file sizes
func WalkFolder(
	path string,
	readDir ReadDir,
	ignoreFunction ShouldIgnoreFolder,
	progress chan<- int,
) *File {
	var wg sync.WaitGroup
	c := make(chan bool, 2*runtime.NumCPU())
	root := walkSubFolderConcurrently(path, 0, nil, ignoringReadDir(ignoreFunction, readDir), c, &wg, progress)
	wg.Wait()

	root.UpdateSize(-1)
	close(progress)
	return root
}

func walkSubFolderConcurrently(
	path string,
	level int,
	parent *File,
	readDir ReadDir,
	c chan bool,
	wg *sync.WaitGroup,
	progress chan<- int,
) *File {
	result := &File{}
	entries := mylog.Check2(readDir(path))

	dirName, name := filepath.Split(path)
	result.Files = make([]*File, 0, len(entries))
	numSubFolders := 0
	defer updateProgress(progress, &numSubFolders)
	var mutex sync.Mutex
	for _, entry := range entries {
		if entry.IsDir() {
			numSubFolders++
			subFolderPath := filepath.Join(path, entry.Name())
			wg.Add(1)
			go func() {
				c <- true
				subFolder := walkSubFolderConcurrently(subFolderPath, level+1, result, readDir, c, wg, progress)
				if subFolder != nil { // Do not include folders that returned error
					mutex.Lock()
					result.Files = append(result.Files, subFolder)
					mutex.Unlock()
				}
				<-c
				wg.Done()
			}()
		} else {
			size := entry.Size()
			file := &File{
				Name:        entry.Name(),
				FullPath:    filepath.Join(path, entry.Name()),
				Size:        size,
				IsDir:       false,
				Depth:       level,
				NumChildren: 0,
				Files:       []*File{},
			}
			mutex.Lock()
			result.Files = append(result.Files, file)
			mutex.Unlock()
		}
	}

	result.FullPath = path
	result.Depth = level
	result.IsDir = true
	if parent != nil {
		result.Name = name
	} else {
		// Root dir
		// TODO unit test this Join
		result.Name = filepath.Join(dirName, name)
	}

	return result
}

func updateProgress(progress chan<- int, count *int) {
	if *count > 0 {
		progress <- *count
	}
}
