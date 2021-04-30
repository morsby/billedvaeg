package pdf

import (
	"strings"

	billedvaeg "github.com/morsby/billedvaeg/api"
	"github.com/phpdave11/gofpdf"
)

var margin = 20.0
var cellspacing = 5.0

func New() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.SetFont("Arial", "", 11)
	return pdf
}

func AddPeople(pdf *gofpdf.Fpdf, people billedvaeg.PersonList, cols int) {
	rows := 3
	pplPerPage := rows * cols
	var cellWidth = (210-margin*2)/float64(cols) - cellspacing

	for n, p := range people {
		if n%pplPerPage == 0 {
			pdf.AddPage()
		}

		x := margin + (cellWidth+cellspacing*2)*float64(n%cols)

		pdf.SetX(x)

		y := margin + float64(n%pplPerPage/cols*85)
		pdf.SetY(y)

		img := pdf.RegisterImageOptions(p.Img, gofpdf.ImageOptions{})
		w := img.Width() / (img.Height() / 60)

		pdf.ImageOptions(p.Img, x+(cellWidth-w)/2, y, 0, 60, false, gofpdf.ImageOptions{}, 0, "")
		pdf.Ln(60)

		addTextLine(pdf, p.Name, x, cellWidth)
		addTextLine(pdf, p.Position.Title, x, cellWidth)
		addTextLine(pdf, p.Suppl, x, cellWidth)
	}
}

func addTextLine(pdf *gofpdf.Fpdf, text string, x float64, cellWidth float64) {
	toUTF8 := pdf.UnicodeTranslatorFromDescriptor("")

	// weird encoding fixes...
	text = strings.Replace(text, "ü", "ü", -1)
	text = strings.Replace(text, "å", "å", -1)

	pdf.SetX(x)
	pdf.CellFormat(cellWidth, 6, toUTF8(text), "0", 0, "C", false, 0, "")
	pdf.Ln(-1)
}
