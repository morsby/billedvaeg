package billedvaeg

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

// cropImage opens an image, crops it to 3x4 format and
// saves it in the same location with the suffix and
// extension "_cropped.jpg" in the same folder.
func cropImage(buf *bytes.Buffer) (*bytes.Buffer, error) {
	img, _, err := image.Decode(buf)
	if err != nil {
		return nil, err
	}

	cImg, err := cutter.Crop(img, cutter.Config{
		Width:   3,
		Height:  4,
		Mode:    cutter.Centered,
		Options: cutter.Ratio,
	})

	cImg = resize.Resize(15*60, 0, cImg, resize.Lanczos3)

	if err != nil {
		return nil, err
	}

	outBuf := bytes.Buffer{}
	jpeg.Encode(&outBuf, cImg, &jpeg.Options{Quality: 75})
	return &outBuf, nil
}
