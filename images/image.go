package images

import (
	"image"
	"image/jpeg"
	"os"
	"path"
	"strings"

	"github.com/oliamb/cutter"
)

func CropImage(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return "", err
	}

	cImg, err := cutter.Crop(img, cutter.Config{
		Width:   3,
		Height:  4,
		Mode:    cutter.Centered,
		Options: cutter.Ratio,
	})

	if err != nil {
		return "", err
	}
	basepath := strings.TrimSuffix(filepath, path.Ext(filepath))
	newFile := basepath + "_cropped.jpg"
	out, _ := os.Create(newFile)
	defer out.Close()
	jpeg.Encode(out, cImg, &jpeg.Options{Quality: 100})
	return newFile, nil
}
