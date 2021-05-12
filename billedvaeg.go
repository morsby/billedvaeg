/* Package billedvaeg contains an api to allow for creation of a PDF document
** of A4-sized pages, containing a number of images of Doctors, their names,
** positions and some supplementary information
 */
package billedvaeg

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
)

// Person contains information on a doctor
type Person struct {
	Name     string
	Position *Position
	Suppl    string
	Img      bytes.Buffer
}

// PersonList contains a slice of *Person
type PersonList []*Person

// Sort sorts a PersonList by positions - and if they're equal,
// by name
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

// Position contains information on a position
type Position struct {
	Title string
	Abbr  string
	Value int
}

// PositionList contains a slice of *Position
type PositionList []*Position

// Positions contains the positions, unmarshalled from `positions.json`
var Positions PositionList

//go:embed embeds/positions.json
var positionsJson []byte

func init() {
	json.Unmarshal(positionsJson, &Positions)
	for n, p := range Positions {
		p.Value = n
	}
}

// ToMap converts a PositionstList (a slice) into a map,
// where the key is the position's 'abbr' (abbreviation).
func (pl PositionList) ToMap() map[string]*Position {
	m := make(map[string]*Position)
	for _, p := range pl {
		m[p.Abbr] = p
	}
	return m
}

// parseMultiformData parses a HTTP MultiFormData request, creating a PersonList
// from it's information and returning the data.
func parseMultiformData(r *http.Request) (*PersonList, error) {
	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called
	files := r.MultipartForm.File["file"]
	updatedPpl := PersonList{}
	var positions = Positions.ToMap()
	for _, file := range files {
		// Open the file
		uploadedFile, err := file.Open()
		if err != nil {
			return nil, err
		}
		// Detect content type:
		buff := make([]byte, 512)
		_, err = uploadedFile.Read(buff)
		if err != nil {
			return nil, err
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			return nil, errors.New("provided file format is not allowed")
		}

		// Reset read position
		_, err = uploadedFile.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}

		// Create a copy
		img := bytes.Buffer{}
		io.Copy(&img, uploadedFile)
		uploadedFile.Close()

		name := r.Form[file.Filename+"-name"][0]
		position := positions[r.Form[file.Filename+"-position"][0]]
		suppl := r.Form[file.Filename+"-suppl"][0]
		// If not a specialist position:
		if position.Value >= Positions.ToMap()["HU"].Value {
			suppl = fmt.Sprintf("Vejleder: %s", suppl)
		}

		updatedPpl = append(updatedPpl, &Person{
			Name:     name,
			Position: position,
			Suppl:    suppl,
			Img:      img,
		})

	}
	return &updatedPpl, nil
}
