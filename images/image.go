package images

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"strings"

	"github.com/morsby/billedvaeg"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

// cropImage opens an image, crops it to 3x4 format and
// saves it in the same location with the suffix and
// extension "_cropped.jpg" in the same folder.
func cropImage(filepath string) (string, error) {
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

	cImg = resize.Resize(15*60, 0, cImg, resize.Lanczos3)

	if err != nil {
		return "", err
	}
	basepath := strings.TrimSuffix(filepath, path.Ext(filepath))
	newFile := basepath + "_cropped.jpg"
	out, _ := os.Create(newFile)
	defer out.Close()
	jpeg.Encode(out, cImg, &jpeg.Options{Quality: 75})
	return newFile, nil
}

// ReadDir reads a dir and takes all images in it, converts them to a
// []Person
func ReadDir(dir string, special bool) (*billedvaeg.PersonList, error) {
	var ppl billedvaeg.PersonList
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !(strings.Contains(file.Name(), ".jpg") ||
			strings.Contains(file.Name(), ".jpeg") ||
			strings.Contains(file.Name(), ".png")) ||
			strings.Contains(file.Name(), "_cropped.jpg") {
			continue
		}
		basepath := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))
		data := strings.Split(basepath, "_")
		if len(data) != 3 {
			return nil, errors.New("Unable to parse filename " + file.Name())
		}
		imgPath, err := cropImage(dir + "/" + file.Name())
		if err != nil {
			return nil, err
		}

		suppl := ""
		if special {
			suppl = strings.Replace(data[2], "-", "/", -1)
		} else {
			suppl = fmt.Sprintf("Vejleder: %s", data[2])
		}

		person := billedvaeg.Person{
			Name:     data[0],
			Position: billedvaeg.Positions[data[1]],
			Suppl:    suppl,
			Img:      imgPath,
		}

		ppl = append(ppl, &person)
	}

	return &ppl, nil
}

// RemoveTmpFiles deletes all created cropped images
// (based on "_cropped.jpg" filename) in the given dir.
func RemoveTmpFiles(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !strings.Contains(file.Name(), "_cropped.jpg") {
			continue
		}
		err := os.Remove(dir + "/" + file.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveTmpDir deletes the given directory and all contents
func RemoveTmpDir(dir string) error {
	return os.RemoveAll(dir)
}
