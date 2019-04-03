package assets

import (
	"bytes"
	"github.com/lxn/walk"
	"image"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
)

func EnsureManifestFile() {

	fileName := filepath.Base(os.Args[0])
	fileName = filepath.Join(filepath.Dir(os.Args[0]), fileName+".exe.manifest")
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		return
	}

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	b, err := Asset("assets/exe.manifest")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := file.Write(b); err != nil {
		log.Fatal(err)
	}

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}

var (
	ImgPinOn     = Image("assets/png16/pin_on.png")
	ImgPinOff    = Image("assets/png16/pin_off.png")
	ImgForward   = Image("assets/png16/forward.png")
	ImgError     = Image("assets/png16/error.png")
	ImgCheckMark = Image("assets/png16/checkmark.png")
)

func Image(path string) walk.Image {

	b, err := Asset(path)
	if err != nil {
		log.Fatalln(err, path)
	}

	x, s, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatalln(err, s, path)
	}
	r, err := walk.NewBitmapFromImage(x)
	if err != nil {
		log.Fatalln(err, s, path)
	}
	return r

}
