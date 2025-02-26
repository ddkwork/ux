package ux

import (
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/stream"
)

func TestName(t *testing.T) {
	g := stream.NewGeneratedFile()
	g.P(`
package ux

import (
	"embed"
	"github.com/ddkwork/golibrary/stream"
	"strings"
)

func svgCallback(value []byte) []byte {
	return []byte(strings.Replace(string(value), "<path ", "<path fill=\"white\"", 1))
}

// 取色
// https://products.eptimize.app/zh/color-convert/rgb-to-rgba
//
//go:embed resources/images/*.svg
var images embed.FS
var (
	svgEmbedFileMap = stream.ReadEmbedFileMap(images, "resources/images")
`)

	for k := range svgEmbedFileMap.Range() {
		// mylog.Info(k, v)
		name := "SvgIcon" + stream.ToCamelUpper(stream.TrimExtension(k))
		g.P(name, "=svgEmbedFileMap.GetMustCallback(", strconv.Quote(k), ",svgCallback)")
	}
	g.P(")")
	stream.WriteGoFile("svgIcons_gen.go", g.Bytes())
}
