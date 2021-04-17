package main

import (
	"github.com/morsby/billedvaeg/images"
	"github.com/morsby/billedvaeg/pdf"
)

func main() {
	imgFolder := "./data"
	ppl, err := images.ReadDir(imgFolder)
	defer images.RemoveTmpFiles(imgFolder)

	if err != nil {
		panic(err)
	}

	doc := pdf.New()
	pdf.AddPeople(doc, ppl)

	err = doc.OutputFileAndClose("Billedvæg.pdf")
	if err != nil {
		panic(err)
	}
}
