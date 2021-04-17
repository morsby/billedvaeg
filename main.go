package main

import (
	"github.com/morsby/billedvaeg/pdf"
)

func main() {

	ppl := []pdf.Person{
		{Name: "Sigurd Morsby Larsen", Position: "I-læge", Img: "500x400.jpg"},
		{Name: "Johanne Kassow", Position: "I-læge", Img: "500x600.jpg"},
		{Name: "Lasse Slumstrup", Position: "I-læge", Img: "500x700.jpg"},
		{Name: "Morten Stokholm", Position: "I-læge", Img: "500x800.jpg"},
		{Name: "Mette Foldager", Position: "I-læge", Img: "500x700.jpg"},
		{Name: "Jenny-Ann Phan", Position: "I-læge", Img: "500x600.jpg"},
		{Name: "Thomas Harbo", Position: "I-læge", Img: "500x400.jpg"},
	}

	doc := pdf.New()
	pdf.AddPeople(doc, ppl)

	err := doc.OutputFileAndClose("Billedvæg.pdf")
	if err != nil {
		panic(err)
	}
}
