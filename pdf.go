package billedvaeg

import (
	"fmt"
	"strings"

	"github.com/phpdave11/gofpdf"
)

var margin = 20.0
var cellspacing = 5.0

type Document struct {
	PDF    *gofpdf.Fpdf
	People []*Person
	Cols   int
	Rows   int
}

func New() Document {
	doc := Document{}
	doc.PDF = gofpdf.New("P", "mm", "A4", "")
	// default values
	doc.Cols = 3
	doc.Rows = 3

	doc.PDF.SetFont("Arial", "", 11)
	return doc
}

func (doc Document) AddPeople(people []*Person) {
	pplPerPage := doc.Rows * doc.Cols
	var cellWidth = (210-margin*2)/float64(doc.Cols) - cellspacing

	for n, p := range people {
		if n%pplPerPage == 0 {
			doc.PDF.AddPage()
		}

		x := margin + (cellWidth+cellspacing*2)*float64(n%doc.Cols)

		doc.PDF.SetX(x)

		y := margin + float64(n%pplPerPage/doc.Cols*85)
		if p.Img != nil {
			doc.PDF.SetY(y)
			imgName := fmt.Sprintf("img-%s", strings.ReplaceAll(p.Name, " ", "-"))
			cImg, err := cropImage(p.Img)
			if err != nil {
				panic(err)
			}
			img := doc.PDF.RegisterImageOptionsReader(imgName, gofpdf.ImageOptions{ImageType: "JPEG"}, cImg)
			w := img.Width() / (img.Height() / 60)

			doc.PDF.ImageOptions(imgName, x+(cellWidth-w)/2, y, 0, 60, false, gofpdf.ImageOptions{}, 0, "")
		}
		doc.PDF.Ln(60)

		doc.addTextLine(p.Name, x, cellWidth)
		doc.addTextLine(p.Position.Title, x, cellWidth)
		doc.addTextLine(p.Suppl, x, cellWidth)
	}
}

func (doc Document) addTextLine(text string, x float64, cellWidth float64) {
	toUTF8 := doc.PDF.UnicodeTranslatorFromDescriptor("")

	// weird encoding fixes...
	text = strings.Replace(text, "ü", "ü", -1)
	text = strings.Replace(text, "å", "å", -1)

	doc.PDF.SetX(x)
	doc.PDF.CellFormat(cellWidth, 6, toUTF8(text), "0", 0, "C", false, 0, "")
	doc.PDF.Ln(-1)
}
