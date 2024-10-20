package guiutils

import (
	"embed"
	"fmt"
	"image"
	"path/filepath"
	"time"

	"gioui.org/text"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux"
	"github.com/ddkwork/ux/filetree/files"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/dustin/go-humanize"
)

//go:embed images/gocleasy-logo.png
var imageFile embed.FS

type State string

const (
	HomeS         State = "home"          // Show the scan button
	LoadingFilesS State = "loadingFilesS" // Show the files to be selected
	SelFilesS     State = "selFileS"      // Show the files to be selected
	DelFilesS     State = "delFileS"      // Show the selected files to be deleted
)

type AppLogic struct {
	theme      *material.Theme   // Store the them of the application
	Files      *files.File       // Used to store the files with their structure
	Selfiles   []*files.File     // Used to store the files that has been selected
	Files2Show []*files.FileShow // Used to store the filest that are going to be rendered
	Appstate   State
}

type (
	C = layout.Context
	D = layout.Dimensions
)

// NewAppLogic Create an instance of AppLogic
func NewAppLogic() *AppLogic {
	return &AppLogic{
		theme:    ux.ThemeDefault().Theme,
		Appstate: HomeS,
	}
}

//	th := material.NewTheme()
//	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

func (a *AppLogic) ReportProgress(win *app.Window, total *int, progress <-chan int) {
	// Controls how frequently to update the application
	const interval = 250 * time.Millisecond
	ticker := time.NewTicker(interval)

	for {
		select {
		case loadedFile, ok := <-progress:

			if ok && (a.Appstate == LoadingFilesS) {
				*total += loadedFile
			} else if !ok {
				a.Appstate = SelFilesS
				win.Invalidate()
				return
			}

		case <-ticker.C:
			win.Invalidate()
		}
	}
}

func showGoCleasyLogo(gtx C, margins layout.Inset) layout.FlexChild {
	// Show logo
	return layout.Flexed(1, func(gtx C) D {
		// Open the file using the file path
		file := mylog.Check2(imageFile.Open("images/gocleasy-logo.png"))

		defer file.Close()

		// Pass the file to the Decode function
		img, format := mylog.Check3(image.Decode(file))

		return margins.Layout(gtx, func(gtx C) D {
			return widget.Image{
				Src: paint.NewImageOp(img),
				Fit: widget.Contain,
			}.Layout(gtx)
		})
	})
}

func (a *AppLogic) HomePage(gtx C, scanbutton *widget.Clickable, initialpathinput *widget.Editor, numfilesdeleted int64, sizeliberated int64) D {
	margins := layout.Inset{
		Top:    unit.Dp(25),
		Bottom: unit.Dp(25),
		Right:  unit.Dp(25),
		Left:   unit.Dp(25),
	}

	var deletedfilesoutput string
	if numfilesdeleted == 0 || sizeliberated == 0 {
		deletedfilesoutput = ""
	} else {
		deletedfilesoutput = fmt.Sprintf("   Deleted %s files and %s", humanize.Comma(numfilesdeleted), humanize.Bytes(uint64(sizeliberated)))
	}

	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceEnd,
	}.Layout(gtx,
		showGoCleasyLogo(gtx, margins),
		layout.Rigid(func(gtx C) D {
			return material.Body1(a.theme, deletedfilesoutput).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Right: unit.Dp(25),
				Left:  unit.Dp(25),
			}.Layout(gtx, func(gtx C) D {
				return material.Editor(a.theme, initialpathinput, " Introduce Initial Path. Leave blank for root path.").Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return margins.Layout(gtx, func(gtx C) D {
				return material.Button(a.theme, scanbutton, "Scan Files").Layout(gtx)
			})
		}),
	)
}

func createTextNLoading(gtx C, th *material.Theme, text string) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		return layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(gtx,
			layout.Rigid(
				layout.Spacer{Width: unit.Dp(25)}.Layout,
			),
			layout.Rigid(func(gtx C) D {
				return material.Body1(th, fmt.Sprintf("Loading file \"%s\"...", text)).Layout(gtx)
			}),
			layout.Rigid(
				layout.Spacer{Width: unit.Dp(25)}.Layout,
			),
			layout.Rigid(func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					return material.Loader(th).Layout(gtx)
				})
			}))
	})
}

func (a *AppLogic) FillFirstLayer2Show() {
	for _, file := range a.Files.Files {
		a.Files2Show = append(a.Files2Show, &files.FileShow{
			File:         file,
			IsSelected:   widget.Bool{},
			ActionButton: widget.Bool{},
		})
	}
}

func (a *AppLogic) ShowLoadingPage(gtx C, actualFilesRead int) D {
	margins := layout.Inset{
		Top:    unit.Dp(25),
		Bottom: unit.Dp(25),
		Right:  unit.Dp(35),
		Left:   unit.Dp(35),
	}

	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceEnd,
	}.Layout(gtx,
		// Space on the top of the window
		layout.Rigid(
			layout.Spacer{Height: unit.Dp(25)}.Layout,
		),
		showGoCleasyLogo(gtx, margins),
		// Show Reading files and loading circle
		createTextNLoading(gtx, a.theme, fmt.Sprintf("%d", actualFilesRead)),
		layout.Rigid(
			layout.Spacer{Height: unit.Dp(25)}.Layout,
		),
	)
}

func selectFilesTableRow(th *material.Theme, file *files.FileShow, numchildren string, filepath string, rowIndex int) []layout.FlexChild {
	return []layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = 1200
			gtx.Constraints.Max.X = 1200
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(layout.Spacer{Width: unit.Dp(file.File.Depth * 25)}.Layout),
				layout.Rigid(func(gtx C) D {
					ux.DrawColumnDivider(gtx, 0)
					return material.CheckBox(th, &file.IsSelected, "").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					ux.DrawColumnDivider(gtx, 1)
					return material.CheckBox(th, &file.ActionButton, filepath).Layout(gtx)
				}),
			)
		}),

		// Ocupy the space in between buttons and text (checkbox and filenames, size and numfiles)
		// layout.Flexed(1, layout.Spacer{}.Layout),
		// Num of files inside the directory (0 if it is a file)
		layout.Rigid(func(gtx C) D {
			ux.DrawColumnDivider(gtx, 2)
			bgColor := ux.RowColor(rowIndex)
			return ux.Background{Color: bgColor}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = 200
				gtx.Constraints.Max.X = 200
				body1 := material.Body1(th, numchildren)
				body1.Alignment = text.Middle
				body1.MaxLines = 1
				body1.Truncator = "..."
				return body1.Layout(gtx)
			})
		}),
		// layout.Rigid(layout.Spacer{Width: unit.Dp(25)}.Layout),
		// Size of the file
		layout.Rigid(func(gtx C) D {
			ux.DrawColumnDivider(gtx, 3)
			bgColor := ux.RowColor(rowIndex)
			return ux.Background{Color: bgColor}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = 200
				gtx.Constraints.Max.X = 200
				body1 := material.Body1(th, humanize.Bytes(uint64(file.File.Size)))
				body1.Alignment = text.Middle
				body1.MaxLines = 1
				body1.Truncator = "..."
				return body1.Layout(gtx)
			})
		}),
	}
}

func drawTableHeader(gtx C, th *material.Theme) D {
	return layout.Flex{
		Axis: layout.Horizontal,
		// Alignment: layout.Middle,
		// Spacing:   layout.SpaceStart,
	}.Layout(gtx,
		// layout.Rigid(layout.Spacer{Width: unit.Dp(75)}.Layout),
		// Name of the file
		layout.Rigid(func(gtx C) D {
			ux.DrawColumnDivider(gtx, 0)
			gtx.Constraints.Min.X = 1200
			gtx.Constraints.Max.X = 1200
			body1 := material.Body1(th, "Path")
			body1.Alignment = text.Middle
			body1.MaxLines = 1
			body1.Truncator = "..."
			return body1.Layout(gtx)
		}),
		// Ocupy the space in between buttons and text (checkbox and filenames, size and numfiles)
		// Num of files inside the directory (0 if it is a file)
		layout.Rigid(func(gtx C) D {
			ux.DrawColumnDivider(gtx, 1)
			gtx.Constraints.Min.X = 200
			gtx.Constraints.Max.X = 200
			body1 := material.Body1(th, "Num Children")
			body1.Alignment = text.Middle
			body1.MaxLines = 1
			body1.Truncator = "..."
			return body1.Layout(gtx)
		}),
		// Size of the file
		layout.Rigid(func(gtx C) D {
			ux.DrawColumnDivider(gtx, 2)
			gtx.Constraints.Min.X = 200
			gtx.Constraints.Max.X = 200
			body1 := material.Body1(th, "Size")
			body1.Alignment = text.Middle
			body1.MaxLines = 1
			body1.Truncator = "..."
			return body1.Layout(gtx)
		}),
	)
}

func deleteFilesTableRow(gtx C, th *material.Theme, field1 string, field2 string, field3 string) D {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(layout.Spacer{Width: unit.Dp(25)}.Layout),
		// Name of the file
		layout.Rigid(func(gtx C) D {
			return material.Body1(th, field1).Layout(gtx)
		}),
		// Ocupy the space in between buttons and text (checkbox and filenames, size and numfiles)
		layout.Flexed(1, layout.Spacer{}.Layout),
		// Num of files inside the directory (0 if it is a file)
		layout.Rigid(func(gtx C) D {
			return material.Body1(th, field2).Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(25)}.Layout),
		// Size of the file
		layout.Rigid(func(gtx C) D {
			return material.Body1(th, field3).Layout(gtx)
		}),
	)
}

func (a *AppLogic) DrawTreeTable(gtx C, nextbutton *widget.Clickable, filelist *widget.List) D {
	widgets := []layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			return drawTableHeader(gtx, a.theme)
		}),
		layout.Rigid(func(gtx C) D {
			return a.drawTreeTableRows(gtx, filelist, "")
		}),
	}

	widgets = append(widgets,
		// Button to confirm selected files
		layout.Rigid(func(gtx C) D {
			margins := layout.Inset{
				Top:    unit.Dp(25),
				Bottom: unit.Dp(25),
				Right:  unit.Dp(35),
				Left:   unit.Dp(35),
			}
			return margins.Layout(gtx, func(gtx C) D {
				return material.Button(a.theme, nextbutton, "Next").Layout(gtx) // todo bug
			})
		}))

	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceStart,
	}.Layout(gtx, widgets...)
}

// Checks if ref is inside slf
func isFileSelected(ref *files.File, slf []*files.File) bool {
	for _, file := range slf {
		if file == ref {
			return true
		}
	}
	return false
}

// Calculates how many files need to be deleted from the slice when a folder is closed
// It loops from the actual position till the end of the slice or till a file with lower
// level than the folder being closed
func getNumFiles2NotShow(pos int, level int, sl []*files.FileShow) int {
	res := 0

	for ; pos < len(sl); pos++ {
		file := sl[pos]
		if file.File.Depth <= level {
			return res
		}
		res++
	}

	return res
}

// It loops over Files2Show and checks if there is any checkbox has been clicked to open a folder.
// It also checks if any folder/file has been selected and adds it to Selfiles
func (a *AppLogic) getFiles2Show(gtx C) {
	var file *files.FileShow

	index := 0
	for index < len(a.Files2Show) {

		file = a.Files2Show[index]

		// Check Open/Close folders
		if file.ActionButton.Update(gtx) && file.File.IsDir {
			if file.ActionButton.Value {
				// Add children from Files2Show (Open folder)

				// Create temporal slice to add to children (Files2Show)
				var slice2add []*files.FileShow
				for _, file2append := range file.File.Files {
					// Check if the file was selected before to add it selected
					slice2add = append(slice2add, &files.FileShow{
						File:         file2append,
						IsSelected:   widget.Bool{Value: isFileSelected(file2append, a.Selfiles)},
						ActionButton: widget.Bool{},
					})
				}

				// Insert temporal slice into files to show
				a.Files2Show = append(a.Files2Show[:index+1], append(slice2add, a.Files2Show[index+1:]...)...)

			} else {
				// Delete children from Files2Show (Close folder)
				numFiles2NotShow := getNumFiles2NotShow(index+1, file.File.Depth, a.Files2Show)
				a.Files2Show = append(a.Files2Show[:index+1], a.Files2Show[index+1+numFiles2NotShow:]...)
			}
		}

		// Check selected files
		if file.IsSelected.Update(gtx) {
			if file.IsSelected.Value {
				// Add file to Selfiles
				a.Selfiles = append(a.Selfiles, file.File)
			} else {
				// Delete file from Selfiles
				for id, delfile := range a.Selfiles {
					if delfile == file.File {
						a.Selfiles = append(a.Selfiles[:id], a.Selfiles[id+1:]...)
					}
				}
			}
		}

		index++

	}
}

// Contains the file Tree
func (a *AppLogic) drawTreeTableRows(gtx C, filelist *widget.List, path string) D {
	// empty the files to show
	numfiles := 0
	a.getFiles2Show(gtx)
	numfiles = len(a.Files2Show)

	if numfiles == 0 {
		return D{}
	}

	return material.List(a.theme, filelist).Layout(gtx, len(a.Files2Show), func(gtx layout.Context, index int) layout.Dimensions {
		file := a.Files2Show[index]
		var widgets []layout.FlexChild

		if file.File.IsDir {
			widgets = selectFilesTableRow(a.theme, file, humanize.Comma(file.File.NumChildren), fmt.Sprintf("%s/", filepath.Join(path, file.File.Name)), index)
		} else {
			widgets = selectFilesTableRow(a.theme, file, "-", filepath.Join(path, file.File.Name), index)
		}
		return layout.Flex{ // row
			Axis:      layout.Horizontal,
			Spacing:   0,
			Alignment: 0,
			WeightSum: 0,
		}.Layout(gtx, widgets...)
	})
}

func (a *AppLogic) ShowDeletingPage(gtx C, comebackbutton *widget.Clickable, copy2clipboard *widget.Clickable, deletebutton *widget.Clickable, filedeletelist *widget.List) D {
	margins := layout.Inset{
		Top:    unit.Dp(15),
		Bottom: unit.Dp(15),
		Right:  unit.Dp(15),
		Left:   unit.Dp(15),
	}

	var totFiles, totSize int64 = 0, 0
	for _, file := range a.Selfiles {
		if file.IsDir {
			totFiles += file.NumChildren
		} else {
			totFiles++
		}
		totSize += file.Size

	}

	return layout.Flex{
		Alignment: layout.Middle,
		Axis:      layout.Vertical,
	}.Layout(gtx,
		// Space on the top of the window
		layout.Rigid(
			layout.Spacer{Height: unit.Dp(25)}.Layout,
		),
		layout.Rigid(func(gtx C) D {
			return deleteFilesTableRow(gtx, a.theme, "Path", "Num Children", "Size")
		}),
		// Show selected files
		layout.Flexed(1, func(gtx C) D {
			return a.selectedFiles(gtx, filedeletelist)
		}),
		// Show total
		layout.Rigid(func(gtx C) D {
			return deleteFilesTableRow(gtx, a.theme, "Total", humanize.Comma(totFiles), humanize.Bytes(uint64(totSize)))
		}),
		// Show control buttons
		layout.Rigid(func(gtx C) D {
			return layout.Flex{
				Alignment: layout.Middle,
				Axis:      layout.Horizontal,
			}.Layout(gtx,
				// Show comeback button
				layout.Flexed(1, func(gtx C) D {
					return margins.Layout(gtx, func(gtx C) D {
						return material.Button(a.theme, comebackbutton, "Back").Layout(gtx)
					})
				}),
				// Show copy to clipboard button
				layout.Flexed(1, func(gtx C) D {
					return margins.Layout(gtx, func(gtx C) D {
						return material.Button(a.theme, copy2clipboard, "Copy to Clipboard").Layout(gtx)
					})
				}),
				// Show delete button
				layout.Flexed(1, func(gtx C) D {
					return margins.Layout(gtx, func(gtx C) D {
						return material.Button(a.theme, deletebutton, "Delete").Layout(gtx)
					})
				}),
			)
		}),
	)
}

func (a *AppLogic) selectedFiles(gtx C, l *widget.List) D {
	return l.List.Layout(gtx, len(a.Selfiles), func(gtx C, index int) D {
		f := a.Selfiles[index]
		var numChildren, fullPath string
		if f.IsDir {
			fullPath = fmt.Sprintf("%s/", f.FullPath)
			numChildren = humanize.Comma(f.NumChildren)
		} else {
			numChildren = "-"
			fullPath = f.FullPath
		}
		return deleteFilesTableRow(gtx, a.theme, fullPath, numChildren, humanize.Bytes(uint64(f.Size)))
	})
}
