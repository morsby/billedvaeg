package main

import (
	"github.com/phpdave11/gofpdf"
	"golang.org/x/text/encoding/charmap"
)

func main() {
	pdf := gofpdf.New("P", "mm", "A4", "")

	type person struct {
		name     string
		position string
		img      string
	}
	ppl := []person{
		{"Sigurd Morsby Larsen", "I-læge", "500x700.jpg"},
		{"Johanne Kassow", "I-læge", "500x700.jpg"},
		{"Lasse Slumstrup", "I-læge", "500x700.jpg"},
		{"Morten Stokholm", "I-læge", "500x700.jpg"},
		{"Mette Foldager", "I-læge", "500x700.jpg"},
		{"Jenny-Ann Phan", "I-læge", "500x700.jpg"},
		{"Thomas Harbo", "I-læge", "500x700.jpg"},
	}

	enc := charmap.Windows1252.NewEncoder()

	margin := 20.0
	cellspacing := 5.0

	cellWidth := (210-margin*2)/2 - cellspacing
	y := margin
	// Simple table
	basicTable := func() {
		// Each cell
		for n, p := range ppl {
			if n%6 == 0 {
				pdf.AddPage()
			}

			x := margin
			if n%2 == 1 {
				x += cellWidth + cellspacing*2
			}
			pdf.SetX(x)

			y = margin + float64(n%6/2*80)
			pdf.SetY(y)

			name, err := enc.String(p.name)
			if err != nil {
				panic(err)
			}
			pos, err := enc.String(p.position)
			if err != nil {
				panic(err)
			}

			img := pdf.RegisterImageOptions("data/"+p.img, gofpdf.ImageOptions{})
			w := img.Width() / (img.Height() / 60)

			pdf.ImageOptions("data/"+p.img, x+(cellWidth-w)/2, y, 0, 60, false, gofpdf.ImageOptions{}, 0, "")
			pdf.Ln(60)
			pdf.SetX(x)
			pdf.CellFormat(cellWidth, 6, name, "0", 0, "C", false, 0, "")
			pdf.Ln(-1)
			pdf.SetX(x)
			pdf.CellFormat(cellWidth, 6, pos, "0", 0, "C", false, 0, "")
			pdf.Ln(-1)
		}
	}

	pdf.SetFont("Arial", "", 12)
	basicTable()

	err := pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		panic(err)
	}
}
