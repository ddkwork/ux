package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/ddkwork/ux/internal/giosvg/internal/svgparser"

	"golang.org/x/text/language"

	"golang.org/x/text/cases"
)

var (
	input  string
	output string
	pkg    string
)

func main() {
	flag.StringVar(&input, "i", "", "folder containing svg images or the path of svg file")
	flag.StringVar(&output, "o", "", "file path to save the go code")
	flag.StringVar(&pkg, "pkg", "", "package name")
	flag.Parse()

	if input == "" {
		panic("invalid input")
	}

	var paths []string
	s, err := os.Stat(input)
	if err != nil {
		panic(err)
	}

	if s.IsDir() {
		if paths, err = filepath.Glob(filepath.Join(input, "*.svg")); err != nil {
			panic(err)
		}
	} else {
		paths = []string{input}
	}

	if pkg == "" {
		if output != "" {
			abs, _ := filepath.Abs(output)
			pkg = filepath.Base(filepath.Dir(abs))
			pkg = strings.Replace(pkg, "", "_", -1)
			pkg = strings.Replace(pkg, "-", "_", -1)
		} else {
			pkg = "assets"
		}
		fmt.Printf("package name not defined, use -pkg flag to set a custom name. Using %s as name instead. \n", pkg)
	}

	out := bytes.NewBuffer(nil)
	fmt.Fprintf(out, `// autogenerated by svggen
package %s

import (
	"image"	
	"image/color"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/ddkwork/ux/giosvg"
)

var _, _, _, _, _, _, _, _ = (*f32.Point)(nil), (*op.Ops)(nil), (*clip.Op)(nil), (*paint.PaintOp)(nil), (*giosvg.Vector)(nil), (*color.NRGBA)(nil), (*layout.Dimensions)(nil), (*image.Image)(nil)
`+"\r\n", pkg)

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		svg, err := svgparser.ReadIcon(f)
		if err != nil {
			panic(err)
		}

		name := filepath.Base(path)
		name = strings.Replace(name, "-", " ", -1)
		name = strings.Replace(name, "_", " ", -1)
		name = cases.Title(language.English).String(strings.Replace(name, filepath.Ext(name), "", -1))
		name = strings.Replace(name, " ", "", -1)

		fmt.Fprintf(out, `var Vector%s giosvg.Vector = func(ops *op.Ops, constraints giosvg.Constraints) layout.Dimensions {`+"\r\n", name)

		fmt.Fprintf(out, `var w, h float32`+"\r\n")
		fmt.Fprintf(out, `if constraints.Max != constraints.Min {`+"\r\n")
		if svg.ViewBox.W >= svg.ViewBox.H {
			fmt.Fprintf(out, `
			d := float32(%f)
			if constraints.Max.Y*d > constraints.Max.X {
				w, h = constraints.Max.X, constraints.Max.X/d
			} else {
				w, h = constraints.Max.Y*d, constraints.Max.Y
			}`, svg.ViewBox.W/svg.ViewBox.H)
		} else {
			fmt.Fprintf(out, `
			d := float32(%f)
			if constraints.Max.X*d > constraints.Max.Y {
				w, h = constraints.Max.Y/d, constraints.Max.Y
			} else {
				w, h = constraints.Max.X, constraints.Max.X*d
			}`, svg.ViewBox.H/svg.ViewBox.W)
		}
		fmt.Fprintf(out, `}`+"\r\n")

		fmt.Fprintf(out, `
		if constraints.Min.X > w {
			w = constraints.Min.X
		}
		if constraints.Min.Y > h {
			h = constraints.Min.Y
		}
`)

		fmt.Fprintf(out, `
var (
	size = f32.Point{X: w / %f, Y: h / %f}
	avg = (size.X + size.Y) / 2
	affBase = f32.Affine2D{}.Scale(f32.Point{X: float32(0 - %f), Y: float32(0 - %f)}, size)
	aff = affBase

	end 		clip.PathSpec
	path		clip.Path
	stroke, outline clip.Stack
)`+"\r\n", svg.ViewBox.W, svg.ViewBox.H, svg.ViewBox.X, svg.ViewBox.Y)

		fmt.Fprintf(out, `_, _, _, _, _, _ = avg, aff, end, path, stroke, outline`+"\r\n")

		for _, v := range svg.SVGPaths {
			if v.Style.FillerColor == nil && v.Style.LinerColor == nil {
				continue
			}

			if t := v.Style.Transform; t != svgparser.Identity {
				fmt.Fprintf(out, `aff = affBase.Mul(f32.NewAffine2D(%f, %f, %f, %f, %f, %f))`+"\r\n", t.A, t.C, t.E, t.B, t.D, t.F)
			} else {
				fmt.Fprintf(out, `aff = affBase`+"\r\n")
			}

			fmt.Fprintf(out, "\r\n"+`path = clip.Path{}`+"\r\n")
			fmt.Fprintf(out, `path.Begin(ops)`+"\r\n")

			for _, op := range v.Path {
				switch op := op.(type) {
				case svgparser.OpMoveTo:
					fmt.Fprintf(out, `path.MoveTo(aff.Transform(f32.Point{X: %f, Y: %f}))`+"\r\n", op.X, op.Y)
				case svgparser.OpLineTo:
					fmt.Fprintf(out, `path.LineTo(aff.Transform(f32.Point{X: %f, Y: %f}))`+"\r\n", op.X, op.Y)
				case svgparser.OpQuadTo:
					fmt.Fprintf(out, `path.QuadTo(aff.Transform(f32.Point{X: %f, Y: %f}), aff.Transform(f32.Point{X: %f, Y: %f}))`+"\r\n", op[0].X, op[0].Y, op[1].X, op[1].Y)
				case svgparser.OpCubicTo:
					fmt.Fprintf(out, `path.CubeTo(aff.Transform(f32.Point{X: %f, Y: %f}), aff.Transform(f32.Point{X: %f, Y: %f}), aff.Transform(f32.Point{X: %f, Y: %f}))`+"\r\n", op[0].X, op[0].Y, op[1].X, op[1].Y, op[2].X, op[2].Y)
				case svgparser.OpClose:
					fmt.Fprintf(out, `path.Close()`+"\r\n")
				}
			}

			paint := func(pattern svgparser.Pattern, opacity float64) {
				switch c := pattern.(type) {
				case svgparser.CurrentColor:
					fmt.Fprintf(out, `paint.PaintOp{}.Add(ops)`+"\r\n")
				case svgparser.PlainColor:
					if opacity < 1 {
						c.NRGBA.A = uint8(math.Round(256 * opacity))
					}
					fmt.Fprintf(out, `paint.ColorOp{Color: color.NRGBA{R: %d, G: %d, B: %d, A: %d}}.Add(ops)`+"\r\n", c.NRGBA.R, c.NRGBA.G, c.NRGBA.B, c.NRGBA.A)
					fmt.Fprintf(out, `paint.PaintOp{}.Add(ops)`+"\r\n")
				}
			}

			fmt.Fprintf(out, `end = path.End()`+"\r\n")

			if v.Style.FillerColor != nil {
				fmt.Fprintf(out, `outline = clip.Outline{Path: end}.Op().Push(ops)`+"\r\n")
				paint(v.Style.FillerColor, v.Style.FillOpacity)
				fmt.Fprintf(out, `outline.Pop()`+"\r\n")
			}
			if v.Style.LinerColor != nil {
				fmt.Fprintf(out, `stroke = clip.Stroke{Path: end, Width: %f * avg}.Op().Push(ops)`+"\r\n", v.Style.LineWidth)
				paint(v.Style.LinerColor, v.Style.LineOpacity)
				fmt.Fprintf(out, `stroke.Pop()`+"\r\n")
			}
		}

		fmt.Fprintf(out, `return layout.Dimensions{Size: image.Point{X: int(w), Y: int(h)}}`+"\r\n")
		fmt.Fprintf(out, `}`+"\r\n\r\n")

		f.Close()
	}

	var save io.WriteCloser
	if output == "" {
		save = os.Stdout
	} else {
		if save, err = os.Create(output); err != nil {
			panic(err)
		}
	}
	defer save.Close()

	result, err := format.Source(out.Bytes())
	if err != nil {
		panic(err)
	}

	save.Write(result)
}
