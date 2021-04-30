package billedvaeg

import (
	"encoding/json"
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

const (
	LO = iota
	OLPro
	UAO
	OL
	AL
	HU
	I
	Ps
)

type Position struct {
	Title string
	Abbr  string
	Value int
}

type Positions []*Position

var PositionsJson string = `[
		{ "title": "Ledende overlæge", "abbr": "LO" },
		{ "title": "Overlæge, professor", "abbr": "OL-Pro" },
		{ "title": "Uddannelsesansvarlig overlæge", "abbr": "UAO" },
		{ "title": "Overlæge", "abbr": "OL" },
		{ "title": "Afdelingslæge", "abbr": "AL" },
		{ "title": "HU Neurologi", "abbr": "HU" },
		{ "title": "Introduktionslæge", "abbr": "I" },
		{ "title": "HU Psykiatri", "abbr": "Ps" }
	  ]`

func (ps Positions) FromJSON() Positions {

	var positions Positions

	json.Unmarshal([]byte(PositionsJson), &positions)
	for n, p := range positions {
		p.Value = n
	}
	return positions
}

func (ps Positions) ToMap() map[string]Position {
	m := make(map[string]Position)
	for _, p := range ps {
		m[p.Abbr] = *p
	}
	return m
}
