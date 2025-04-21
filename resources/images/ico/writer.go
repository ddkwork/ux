package ico

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"image"
	"image/draw"
	"image/png"
	"io"

	"github.com/ddkwork/golibrary/mylog"
)

func Encode(w io.Writer, im image.Image) {
	mylog.CheckNil(w)
	b := im.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, im, b.Min, draw.Src)

	header := head{
		0,
		1,
		1,
	}
	entry := direntry{
		Plane:  1,
		Bits:   32,
		Offset: 22,
	}

	pngbuffer := new(bytes.Buffer)
	pngwriter := bufio.NewWriter(pngbuffer)
	mylog.Check(png.Encode(pngwriter, m))
	mylog.Check(pngwriter.Flush())
	entry.Size = uint32(len(pngbuffer.Bytes()))

	bounds := m.Bounds()
	entry.Width = uint8(bounds.Dx())
	entry.Height = uint8(bounds.Dy())
	bb := new(bytes.Buffer)
	mylog.Check(binary.Write(bb, binary.LittleEndian, header))
	mylog.Check(binary.Write(bb, binary.LittleEndian, entry))
	mylog.Check2(w.Write(bb.Bytes()))
	mylog.Check2(w.Write(pngbuffer.Bytes()))
}
