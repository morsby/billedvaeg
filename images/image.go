package images

import (
	"errors"
	"image"
	"image/jpeg"
	"os"
	"path"
	"strings"

	"github.com/morsby/billedvaeg"
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

// ReadDir reads a dir and takes all images in it, converts them to a
// []Person
func ReadDir(dir string) (*billedvaeg.PersonList, error) {
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

		person := billedvaeg.Person{
			Name:     data[0],
			Position: data[1],
			Mentor:   data[2],
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
