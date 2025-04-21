package ico

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
)

func sqDiffUInt8(x, y uint8) uint64 {
	d := uint64(x) - uint64(y)
	return d * d
}

func fastCompare(img1, img2 *image.NRGBA) (int64, error) {
	if img1.Bounds() != img2.Bounds() {
		return 0, fmt.Errorf("image bounds not equal: %+v, %+v", img1.Bounds(), img2.Bounds())
	}

	accumError := int64(0)
	for i := range len(img1.Pix) {
		accumError += int64(sqDiffUInt8(img1.Pix[i], img2.Pix[i]))
	}

	return int64(math.Sqrt(float64(accumError))), nil
}

func TestDecodeConfig(t *testing.T) {
	t.Parallel()
	mylog.Call(func() {
		file := "testdata/golang.ico"
		copyFile := "testdata/golang.png"

		reader := mylog.Check2(os.Open(file))
		icoImage := mylog.Check2(DecodeConfig(reader))
		mylog.Check(reader.Close())

		reader = mylog.Check2(os.Open(copyFile))
		pngImage := mylog.Check2(png.DecodeConfig(reader))
		mylog.Check(reader.Close())
		if icoImage != pngImage {
			t.Errorf("%v - %v", icoImage, pngImage)
		}
	})
}

func TestDecode(t *testing.T) {
	t.Parallel()
	mylog.Call(func() {
		file := "testdata/golang.ico"
		copyFile := "testdata/golang.png"
		reader := mylog.Check2(os.Open(file))
		icoImage := mylog.Check2(Decode(reader))
		if icoImage == nil {
			return
		}
		mylog.Check(reader.Close())

		reader = mylog.Check2(os.Open(copyFile))
		pngImage := mylog.Check2(png.Decode(reader))
		if pngImage == nil {
			return
		}
		mylog.Check(reader.Close())
		if icoImage == nil || !icoImage.Bounds().Eq(pngImage.Bounds()) {
			t.Fatal("bounds differ")
		}
		inrgba := icoImage

		pnrgba := pngImage
		if b := mylog.Check2(fastCompare(inrgba.(*image.NRGBA), pnrgba.(*image.NRGBA))); b != 670 {
			t.Fatalf("pix differ %d \n", b)
		}
	})
}
