package billedvaeg

// this file contains all the code related to the pdf generation.

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/jpeg"
	"strings"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"github.com/phpdave11/gofpdf"
)

var margin = 20.0
var cellspacing = 5.0

// Document contains all necessary information on the people, the positions and
// layout of the document -- which is also contained.
type Document struct {
	PDF       *gofpdf.Fpdf
	People    []*Person
	Positions map[int]*Position
	Cols      int
	Rows      int
}

// New creates a new document with some default values (Cols = 3, Rows = 3).
func New() Document {
	doc := Document{}
	doc.PDF = gofpdf.New("P", "mm", "A4", "")
	// default values
	doc.Cols = 3
	doc.Rows = 3

	doc.PDF.SetFont("Arial", "", 11)
	return doc
}

// the placeholder image if no image is provided
//go:embed placeholder.jpg
var placeholderFile embed.FS

// Generate generates the actual pages; inserting the people sorted.
func (doc Document) Generate(sort bool) error {
	people := doc.People
	if sort {
		// sort people by position > name
		people = SortPersons(doc.People, sort)
	}

	pplPerPage := doc.Rows * doc.Cols
	var cellWidth = (210-margin*2)/float64(doc.Cols) - cellspacing

	for n, p := range people {
		if n%pplPerPage == 0 {
			doc.PDF.AddPage()
		}

		// calculate x position on page
		x := margin + (cellWidth+cellspacing*2)*float64(n%doc.Cols)

		doc.PDF.SetX(x)

		y := margin + float64(n%pplPerPage/doc.Cols*85)
		doc.PDF.SetY(y)

		if p.Img == nil {
			img, err := placeholderFile.Open("placeholder.jpg")
			if err != nil {
				return err
			}
			p.ImageFromReader(img)
		}

		cImg, err := cropImage(p.Img)
		if err != nil {
			return err
		}

		imgName := fmt.Sprintf("img-%02d", n)
		img := doc.PDF.RegisterImageOptionsReader(imgName, gofpdf.ImageOptions{ImageType: "JPEG"}, cImg)
		w := img.Width() / (img.Height() / 60)

		doc.PDF.ImageOptions(imgName, x+(cellWidth-w)/2, y, 0, 60, false, gofpdf.ImageOptions{}, 0, "")

		doc.PDF.Ln(60)

		doc.addTextLine(p.Name, x, cellWidth)
		doc.addTextLine(doc.Positions[p.PositionID].Title, x, cellWidth)
		doc.addTextLine(p.Suppl, x, cellWidth)
	}
	return nil
}

func (doc Document) addTextLine(text string, x float64, cellWidth float64) {
	toUTF8 := doc.PDF.UnicodeTranslatorFromDescriptor("")

	// weird encoding fixes...
	text = strings.ReplaceAll(text, "ü", "ü")
	text = strings.ReplaceAll(text, "å", "å")

	doc.PDF.SetX(x)
	doc.PDF.CellFormat(cellWidth, 6, toUTF8(text), "0", 0, "C", false, 0, "")
	doc.PDF.Ln(-1)
}

// cropImage opens an image, crops it to 3x4 format, encodes it to JPEG and returns it.
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
