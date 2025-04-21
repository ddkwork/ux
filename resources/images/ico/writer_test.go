package ico

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
)

func TestEncode(t *testing.T) {
	t.Parallel()
	mylog.Call(func() {
		origfile := "testdata/golang.ico"
		file := "testdata/golang_test.ico"
		f := mylog.Check2(os.Open("testdata/golang.png"))
		img := mylog.Check2(png.Decode(f))
		mylog.Check(f.Close())

		newFile := mylog.Check2(os.Create(filepath.Join(file)))
		Encode(newFile, img)
		mylog.Check(newFile.Close())

		f = mylog.Check2(os.Open(origfile))
		inrgba := mylog.Check2(Decode(f))
		mylog.Check(f.Close())

		newFile = mylog.Check2(os.Open(file))
		pnrgba := mylog.Check2(Decode(newFile))
		mylog.Check(newFile.Close())

		if b := mylog.Check2(fastCompare(inrgba.(*image.NRGBA), pnrgba.(*image.NRGBA))); b != 0 {
			t.Fatalf("pix differ %d\n", b)
		}
	})
}
