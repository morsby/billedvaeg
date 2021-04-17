package pdf

import (
	"fmt"

	"github.com/morsby/billedvaeg"
	"github.com/phpdave11/gofpdf"
	"golang.org/x/text/encoding/charmap"
)

var margin = 20.0
var cellspacing = 5.0
var cellWidth = (210-margin*2)/2 - cellspacing
var enc = charmap.Windows1252.NewEncoder()

func New() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)
	return pdf
}

func AddPeople(pdf *gofpdf.Fpdf, people []billedvaeg.Person) {
	for n, p := range people {
		if n%6 == 0 {
			pdf.AddPage()
		}

		x := margin
		if n%2 == 1 {
			x += cellWidth + cellspacing*2
		}
		pdf.SetX(x)

		y := margin + float64(n%6/2*85)
		pdf.SetY(y)

		img := pdf.RegisterImageOptions(p.Img, gofpdf.ImageOptions{})
		w := img.Width() / (img.Height() / 60)

		pdf.ImageOptions(p.Img, x+(cellWidth-w)/2, y, 0, 60, false, gofpdf.ImageOptions{}, 0, "")
		pdf.Ln(60)

		addTextLine(pdf, p.Name, x)
		addTextLine(pdf, p.Position, x)
		addTextLine(pdf, fmt.Sprintf("Vejleder: %s", p.Mentor), x)

	}
}

func addTextLine(pdf *gofpdf.Fpdf, text string, x float64) {
	encString, err := enc.String(text)
	if err != nil {
		panic(err)
	}
	pdf.SetX(x)
	pdf.CellFormat(cellWidth, 6, encString, "0", 0, "C", false, 0, "")
	pdf.Ln(-1)
}
