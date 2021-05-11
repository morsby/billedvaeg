package billedvaeg

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

// Person contains information on a doctor
type Person struct {
	Name     string
	Position Position
	Suppl    string
	Img      string
}

type PersonList []*Person

func (ppl PersonList) Sort() {
	sort.Slice(ppl, func(i, j int) bool {
		// Same positions, sort by name instead
		if ppl[i].Position.Value == ppl[j].Position.Value {
			return ppl[i].Name < ppl[j].Name
		}

		// sort by name
		return ppl[i].Position.Value < ppl[j].Position.Value
	})
}

const (
	LO = iota
	OLPro
	UAO
	OL
	AL
	HU
	I
	Ps
)

type Position struct {
	Title string
	Abbr  string
	Value int
}

type Positions []*Position

//go:embed positions.json
var PositionsJson []byte

func (ps Positions) FromJSON() Positions {

	var positions Positions

	json.Unmarshal(PositionsJson, &positions)
	for n, p := range positions {
		p.Value = n
	}
	return positions
}

func (ps Positions) ToMap() map[string]Position {
	m := make(map[string]Position)
	for _, p := range ps {
		m[p.Abbr] = *p
	}
	return m
}

type FormOptions struct {
	MAX_UPLOAD_SIZE int64
	Specialists     bool
}

func HandleMultiformData(r *http.Request, opts FormOptions) (*PersonList, error) {
	/* Upload files */

	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called
	files := r.MultipartForm.File["file"]
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	imgFolder := path.Join(os.TempDir(), timestamp)
	err := os.MkdirAll(imgFolder, os.ModePerm)
	if err != nil {
		return nil, err
	}

	for _, fileHeader := range files {
		// Restrict the size of each uploaded file to 1MB.
		// To prevent the aggregate size from exceeding
		// a specified value, use the http.MaxBytesReader() method
		// before calling ParseMultipartForm()
		if fileHeader.Size > opts.MAX_UPLOAD_SIZE {
			return nil, errors.New("uploaded image is too big: %s")
		}

		// Open the file
		uploadedFile, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}

		buff := make([]byte, 512)
		_, err = uploadedFile.Read(buff)
		if err != nil {
			return nil, err
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			return nil, errors.New("provided file format is not allowed")
		}

		_, err = uploadedFile.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}

		filenameWithLowercaseExt := strings.TrimSuffix(fileHeader.Filename, path.Ext(fileHeader.Filename)) + strings.ToLower(path.Ext(fileHeader.Filename))
		tmpFile, err := os.Create(path.Join(imgFolder, filenameWithLowercaseExt))
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(tmpFile, uploadedFile)
		if err != nil {
			return nil, err
		}

		uploadedFile.Close()
		tmpFile.Close()
	}
	return ReadDir(imgFolder, opts.Specialists, r.MultipartForm.Value)
}
