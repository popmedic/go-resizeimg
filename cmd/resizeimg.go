package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
	"github.com/popmedic/go-logger/log"
)

var (
	width   int
	height  int
	imgPath string
	newPath string
)

func init() {
	flag.IntVar(&width, "w", 0, "width of new image")
	flag.IntVar(&height, "h", 0, "height of new image")
	flag.StringVar(&imgPath, "in", "", "path to image to resize")
	flag.StringVar(&newPath, "out", getDefaultNewFile(), "path to new rezised image file")
}

func main() {
	flag.Parse()

	data, err := ioutil.ReadFile(imgPath)
	if nil != err {
		log.Fatal(os.Exit, err)
	}

	image, _, err := image.Decode(bytes.NewReader(data))
	if nil != err {
		log.Fatal(os.Exit, err)
	}

	newImage := resize.Resize(uint(width), uint(height), image, resize.Lanczos3)

	out, err := os.Create(newPath)
	if nil != err {
		log.Fatal(os.Exit, err)
	}

	switch strings.ToLower(filepath.Ext(imgPath)) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(out, newImage, nil)
	case ".png":
		err = png.Encode(out, newImage)
	default:
		err = errors.New("unknown type \"" + strings.ToUpper(filepath.Ext(imgPath)) + "\"")
	}
	if nil != err {
		log.Fatal(os.Exit, err)
	}
}

func getDefaultNewFile() string {
	return filepath.Join(
		filepath.Dir(imgPath),
		fmt.Sprintf(
			"%s@%dx%d%s",
			strings.TrimSuffix(filepath.Base(imgPath), filepath.Ext(imgPath)),
			width,
			height,
			filepath.Ext(imgPath),
		),
	)
}
