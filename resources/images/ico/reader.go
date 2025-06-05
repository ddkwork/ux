package ico

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"

	"github.com/ddkwork/golibrary/std/mylog"
	"golang.org/x/image/bmp"
)

func init() {
	image.RegisterFormat("ico", "\x00\x00\x01\x00?????\x00", Decode, DecodeConfig)
}

func Decode(r io.Reader) (image.Image, error) {
	var d decoder
	d.decode(r)
	return d.images[0], nil
}

func DecodeAll(r io.Reader) []image.Image {
	var d decoder
	d.decode(r)
	return d.images
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	var (
		d   decoder
		cfg image.Config
		err error
	)
	d.decodeHeader(r)
	d.decodeEntries(r)
	e := d.entries[0]
	buf := make([]byte, e.Size+14)
	n, err := io.ReadFull(r, buf[14:])
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return cfg, err
	}
	buf = buf[:14+n]
	if n > 8 && bytes.Equal(buf[14:22], pngHeader) {
		return png.DecodeConfig(bytes.NewReader(buf[14:]))
	}

	d.forgeBMPHead(buf, &e)
	return bmp.DecodeConfig(bytes.NewReader(buf))
}

type direntry struct {
	Width   byte
	Height  byte
	Palette byte
	_       byte
	Plane   uint16
	Bits    uint16
	Size    uint32
	Offset  uint32
}

type head struct {
	Zero   uint16
	Type   uint16
	Number uint16
}

type decoder struct {
	head    head
	entries []direntry
	images  []image.Image
}

func (d *decoder) decode(r io.Reader) {
	d.decodeHeader(r)
	d.decodeEntries(r)
	d.images = make([]image.Image, d.head.Number)
	for i := range d.entries {
		e := &(d.entries[i])
		data := make([]byte, e.Size+14)
		n := mylog.Check2(io.ReadFull(r, data[14:]))
		data = data[:14+n]
		if n > 8 && bytes.Equal(data[14:22], pngHeader) {
			p := mylog.Check2(png.Decode(bytes.NewReader(data[14:])))
			d.images[i] = p
		} else {
			maskData := d.forgeBMPHead(data, e)
			if maskData != nil {
				data = data[:n+14-len(maskData)]
			}
			b := mylog.Check2(bmp.Decode(bytes.NewReader(data)))
			d.images[i] = b

			bounds := d.images[i].Bounds()
			mask := image.NewAlpha(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
			masked := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
			for row := range int(e.Height) {
				for col := range int(e.Width) {
					if maskData != nil {
						rowSize := (int(e.Width) + 31) / 32 * 4
						if (maskData[row*rowSize+col/8]>>(7-uint(col)%8))&0x01 != 1 {
							mask.SetAlpha(col, int(e.Height)-row-1, color.Alpha{A: 255})
						}
					}
					if e.Bits == 32 {
						rowSize := (int(e.Width)*32 + 31) / 32 * 4
						offset := int(binary.LittleEndian.Uint32(data[10:14]))
						alphaPosition := offset + row*rowSize + col*4 + 3
						var alpha color.Alpha
						if len(data) > alphaPosition {
							alpha = color.Alpha{A: data[alphaPosition]}
						}
						mask.SetAlpha(col, int(e.Height)-row-1, alpha)
					}
				}
			}
			draw.DrawMask(masked, masked.Bounds(), d.images[i], bounds.Min, mask, bounds.Min, draw.Src)
			d.images[i] = masked
		}
	}
}

func (d *decoder) decodeHeader(r io.Reader) {
	mylog.Check(binary.Read(r, binary.LittleEndian, &(d.head)))
	if d.head.Zero != 0 || d.head.Type != 1 {
		mylog.Check(fmt.Errorf("corrupted head: [%x,%x]", d.head.Zero, d.head.Type))
	}
}

func (d *decoder) decodeEntries(r io.Reader) {
	n := int(d.head.Number)
	d.entries = make([]direntry, n)
	for i := range n {
		mylog.Check(binary.Read(r, binary.LittleEndian, &(d.entries[i])))
	}
}

func (d *decoder) forgeBMPHead(buf []byte, e *direntry) (mask []byte) {
	data := buf[14:]
	imageSize := len(data)
	if e.Bits != 32 {
		maskSize := (int(e.Width) + 31) / 32 * 4 * int(e.Height)
		imageSize -= maskSize
		if imageSize <= 0 {
			return
		}
		mask = data[imageSize:]
	}

	copy(buf[0:2], "\x42\x4D")
	dibSize := binary.LittleEndian.Uint32(data[:4])
	w := binary.LittleEndian.Uint32(data[4:8])
	h := binary.LittleEndian.Uint32(data[8:12])
	if h > w {
		binary.LittleEndian.PutUint32(data[8:12], h/2)
	}

	binary.LittleEndian.PutUint32(buf[2:6], uint32(imageSize))

	numColors := binary.LittleEndian.Uint32(data[32:36])
	bits := binary.LittleEndian.Uint16(data[14:16])

	switch bits {
	case 1, 2, 4, 8:
		x := uint32(1 << bits)
		if numColors == 0 || numColors > x {
			numColors = x
		}
	default:
		numColors = 0
	}

	var numColorsSize uint32
	switch dibSize {
	case 12, 64:
		numColorsSize = numColors * 3
	default:
		numColorsSize = numColors * 4
	}
	offset := 14 + dibSize + numColorsSize
	if dibSize > 40 && int(dibSize-4) <= len(data) {
		offset += binary.LittleEndian.Uint32(data[dibSize-8 : dibSize-4])
	}
	binary.LittleEndian.PutUint32(buf[10:14], offset)
	return
}

var pngHeader = []byte{'\x89', 'P', 'N', 'G', '\r', '\n', '\x1a', '\n'}
