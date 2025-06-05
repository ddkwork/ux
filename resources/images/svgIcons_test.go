package images

import (
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/std/stream"
)

func TestName(t *testing.T) {
	g := stream.NewGeneratedFile()
	g.P(`
package images

import (
	"embed"
	"github.com/ddkwork/golibrary/std/stream"
	"strings"
)


func svgCallback(value []byte) []byte {
	if strings.Contains(string(value), "fill=\"none\"") {
		//return []byte(strings.Replace(string(value), "fill=\"none\"", "fill=\"white\"", 1))
	}
	return []byte(strings.Replace(string(value), "<path", "<path fill=\"white\"", 1))
}

// 取色
// https://products.eptimize.app/zh/color-convert/rgb-to-rgba
//
//go:embed images/*.svg
var svgFs embed.FS
var (
	svgEmbedFileMap                = stream.ReadEmbedFileMap(svgFs, "images")
`)

	for k := range svgEmbedFileMap.Range() {
		name := "SvgIcon" + stream.ToCamelUpper(stream.TrimExtension(k))
		g.P(name, "= svgEmbedFileMap.GetMustCallback(", strconv.Quote(k), ",svgCallback)")
	}
	g.P(")")
	stream.WriteGoFile("svgIcons_gen.go", g.Bytes())
}
