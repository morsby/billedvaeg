package billedvaeg

import (
	"sort"
)

// Person contains information on a doctor
type Person struct {
	Name     string
	Position string
	Mentor   string
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

var positions = map[string]int{
	"LO":  LO,
	"UAO": UAO,
	"OL":  OL,
	"AL":  AL,
	"HU":  HU,
	"I":   I,
	"Ps":  Ps,
}

func (ppl PersonList) Sort() {
	sort.Slice(ppl, func(i, j int) bool {
		// Same positions, sort by name instead
		if positions[ppl[i].Position] == positions[ppl[j].Position] {
			return ppl[i].Name < ppl[j].Name
		}

		// sort by name
		return positions[ppl[i].Position] < positions[ppl[j].Position]
	})
}
