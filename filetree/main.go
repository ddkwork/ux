package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux"
	"github.com/ddkwork/ux/filetree/files"
	"github.com/ddkwork/ux/filetree/guiutils"
	"github.com/ddkwork/ux/filetree/ignore"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"golang.design/x/clipboard"
)

var filesFromDirsBeingLoaded = make(chan string, 10) // To send files being scanned inside a directory that has been clicked to be expanded

func calculateDirSize(basePath string) (int64, int64, error) {
	var size int64 = 0
	var numChildren int64 = 0
	mylog.Check(filepath.WalkDir(basePath, func(path string, dire os.DirEntry, err error) error {
		// Get the size if not a directory
		fileinfo := mylog.Check2(os.Stat(path))
		if err == nil {
			size += fileinfo.Size()
			numChildren++
		}

		// Continue even if you cannot read one specific file
		return nil
	}))

	numChildren--

	return size, numChildren, err
}

func DeleteFiles(selectedFiles []*files.File) (int64, int64) {
	var errSlice []error

	var numFiles, sizeLiberated int64

	// Loop over selected files and delete them
	for _, file := range selectedFiles {
		if file.IsDir {
			numFiles += file.NumChildren
		} else {
			numFiles++
		}
		sizeLiberated += file.Size
		log.Print("WARNING: If you are testing, you may want to comment the following lines")
		mylog.Check(os.RemoveAll(file.FullPath))

	}

	for _, er1 := range errSlice {
		// If there is any error with any file/folder you can handle it here
		log.Print(er1)
	}

	return numFiles, sizeLiberated
}

func getRootPath() string {
	switch runtime.GOOS {
	case "windows":
		return getWindowsRootPath()
	case "darwin":
		// macOS
		return "/"
	case "android", "linux":
		// For Android apps you will need to request permission for reading from external folders.
		// I was not able to perform that with gioui or golang, maybe you need to create a connector for JAVA
		return "/"
	case "ios":
		// iOS apps are sandboxed too, so the root path will not be directly accessible
		return getIOSRootPath()
	default:
		return "Unknown OS"
	}
}

func getWindowsRootPath() string {
	// On Windows, the root path is typically the drive where the OS is installed,
	// so we need to get the current drive and concatenate it with the path separator.
	return filepath.VolumeName(os.Getenv("SystemDrive")) + string(filepath.Separator)
}

func getIOSRootPath() string {
	// In iOS, the application's root path is restricted, but you can use other directories like the Documents directory.
	// This is just an example of how you could handle it, but it's not the actual root path.
	documentsDir := mylog.Check2(os.UserHomeDir())

	return filepath.Join(documentsDir, "Documents")
}

// Deletes a string from a slice if exists
func deleteStringFromSlice(str string, slice []string) []string {
	// Find and remove the string from the slide
	for i, v := range slice {
		if v == str {
			slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	return slice
}

func printTree(file *files.File, indent string) {
	fmt.Printf("%s%s\n", indent, file.FullPath)
	for _, subFile := range file.Files {
		printTree(subFile, indent+"\t")
	}
}

// func getSelectedFiles(children []*files.File, selfiles *[]*files.File) []*files.File {
//
// 	// create a pointer to every file
// 	var filep *files.File
//
// 	for index := range children {
// 		filep = children[index]
//
// 		// If it is selected add to selfiles
// 		if filep.IsSelected.Value {
// 			*selfiles = append(*selfiles, filep)
// 		}
//
// 		if (filep.IsDir) && // It is a directory
// 			((filep.Files != nil) && (len(filep.Files) > 0)) && // It contains files inside
// 			!filep.IsSelected.Value { // It is not a selected folder
//
// 			*selfiles = getSelectedFiles(filep.Files, selfiles) // Look for selected files inside
// 		}
// 	}
//
// 	return *selfiles
// }

func copyFilesInClipboard(selfiles []*files.File) {
	result := ""

	for _, file := range selfiles {
		result += "\"" + file.FullPath + "\" "
	}

	clipboard.Write(clipboard.FmtText, []byte(result))
}

func Run(w *app.Window) error {
	applogic := guiutils.NewAppLogic()

	// ops are the operations from the UI
	var ops op.Ops

	// Widget declarations
	var scanButton widget.Clickable
	var deleteButton widget.Clickable
	var comeBackButton widget.Clickable
	var copy2clipboard widget.Clickable
	var nextButton widget.Clickable
	var initialPathInput widget.Editor
	fileList := widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}
	fileDeleteList := widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}

	scanFilesLoadingChan := make(chan int) // Used to transmit how many files have been read
	var numFilesDeleted int64 = 0
	var sizeLiberated int64 = 0

	var initialpath string

	totalFilesReadShow := 0
	mylog. // Used to maintain a count of the files read
		Check(

			// Initialize clipboard so you can set the clipboard
			clipboard.Init())

	// Listen for events in the window
	for {
		switch e := w.Event().(type) {
		// Window closed
		case app.DestroyEvent:
			return e.Err

		// Actions in the window apart from closing
		case app.FrameEvent:

			gtx := app.NewContext(&ops, e)
			ux.BackgroundDark(gtx)

			//
			// ACTIONS TO CHANGE THE STATE OF THE APPLICATION ***
			//
			// Goes from homePage to selection of files
			if scanButton.Clicked(gtx) {

				// reset file directory
				applogic.Files = nil

				initialpath = initialPathInput.Text()
				if initialpath == "" {
					initialpath = getRootPath()
				}

				// Test the introduced path, if not good, use the root path
				_ := mylog.Check2(os.ReadDir(initialpath))

				// If there is no problem, continue

				// Add first level of files to be shown

			}

			// Go to confirm deleting the files
			if nextButton.Clicked(gtx) {
				// applogic.Selfiles = getSelectedFiles(applogic.Files.Files, &applogic.Selfiles)
				applogic.Appstate = guiutils.DelFilesS
			}

			// Go back to selecting the files
			if comeBackButton.Clicked(gtx) {
				applogic.Appstate = guiutils.SelFilesS
			}

			// copy files in clipboard
			if copy2clipboard.Clicked(gtx) {
				copyFilesInClipboard(applogic.Selfiles)
			}

			// Delete the files show a message of number of files deleted and amount of memory freed
			if deleteButton.Clicked(gtx) {
				numFilesDeleted, sizeLiberated = DeleteFiles(applogic.Selfiles)
				applogic.Appstate = guiutils.HomeS
			}
			// ACTIONS TO CHANGE THE STATE OF THE APPLICATION ***

			//
			// STATES OF THE APPLICATION ***
			// What template to render based on applogic state
			//
			switch applogic.Appstate {

			case guiutils.HomeS:
				applogic.HomePage(gtx, &scanButton, &initialPathInput, numFilesDeleted, sizeLiberated)

			case guiutils.LoadingFilesS:
				applogic.ShowLoadingPage(gtx, totalFilesReadShow)

			case guiutils.SelFilesS:
				applogic.DrawTreeTable(gtx, &nextButton, &fileList)

			case guiutils.DelFilesS:
				applogic.ShowDeletingPage(gtx, &comeBackButton, &copy2clipboard, &deleteButton, &fileDeleteList)

			}
			// STATES OF THE APPLICATION ***

			e.Frame(gtx.Ops)

		}
	}
}

func main() {
	go func() {
		// create window
		w := new(app.Window)
		w.Option(
			app.Title("Gocleasy"),
			app.Size(unit.Dp(550), unit.Dp(550)),
		)

		// Run main loop
		if mylog.Check(Run(w)); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
