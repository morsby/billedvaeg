package pdf

import (
	"github.com/morsby/billedvaeg/images"
	"github.com/phpdave11/gofpdf"
	"golang.org/x/text/encoding/charmap"
)

type Person struct {
	Name     string
	Position string
	Img      string
}

var margin = 20.0
var cellspacing = 5.0
var cellWidth = (210-margin*2)/2 - cellspacing
var enc = charmap.Windows1252.NewEncoder()

func New() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)
	return pdf
}

func AddPeople(pdf *gofpdf.Fpdf, people []Person) {
	for n, p := range people {
		if n%6 == 0 {
			pdf.AddPage()
		}

		x := margin
		if n%2 == 1 {
			x += cellWidth + cellspacing*2
		}
		pdf.SetX(x)

		y := margin + float64(n%6/2*80)
		pdf.SetY(y)

		name, err := enc.String(p.Name)
		if err != nil {
			panic(err)
		}
		pos, err := enc.String(p.Position)
		if err != nil {
			panic(err)
		}
		imgPath, err := images.CropImage("data/" + p.Img)
		if err != nil {
			panic(err)
		}
		img := pdf.RegisterImageOptions(imgPath, gofpdf.ImageOptions{})
		w := img.Width() / (img.Height() / 60)

		pdf.ImageOptions(imgPath, x+(cellWidth-w)/2, y, 0, 60, false, gofpdf.ImageOptions{}, 0, "")
		pdf.Ln(60)
		pdf.SetX(x)
		pdf.CellFormat(cellWidth, 6, name, "0", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.SetX(x)
		pdf.CellFormat(cellWidth, 6, pos, "0", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}
}
