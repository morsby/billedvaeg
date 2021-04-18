package billedvaeg

import (
	"sort"
)

// Person contains information on a doctor
type Person struct {
	Name     string
	Position Position
	Suppl    string
	Img      string
}

type PersonList []*Person

const (
	LO = iota
	UAO
	OL
	AL
	HU
	I
	Ps
)

type Position struct {
	Title string
	Value int
}

var Positions = map[string]Position{
	"LO":  {"Ledende overlæge", LO},
	"UAO": {"Uddannelsesansvarlig overlæge", UAO},
	"OL":  {"Overlæge", OL},
	"AL":  {"Afdelingslæge", AL},
	"HU":  {"HU Neurologi", HU},
	"I":   {"Introduktionslæge", I},
	"Ps":  {"HU Psykiatri", Ps},
}

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
