package billedvaeg

import (
	"reflect"
	"testing"
)

func TestSortPersons(t *testing.T) {
	type expect struct {
		input  []*Person
		output []*Person
	}

	pos1nameA := &Person{Name: "A", PositionOrder: 1}
	pos1nameB := &Person{Name: "B", PositionOrder: 1}
	pos2nameA := &Person{Name: "A", PositionOrder: 2}

	ts := []expect{
		{[]*Person{pos1nameA, pos1nameB}, []*Person{pos1nameA, pos1nameB}},
		{[]*Person{pos1nameB, pos1nameA}, []*Person{pos1nameA, pos1nameB}},
		{[]*Person{pos2nameA, pos1nameB}, []*Person{pos1nameB, pos2nameA}},
	}

	for _, v := range ts {
		x := SortPersons(v.input, true)
		if !reflect.DeepEqual(x, v.output) {
			t.Errorf("Input: %v, got: %v; wanted: %v", v.input, x, v.output)
		}
	}
}
