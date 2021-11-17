/* Package billedvaeg contains an api to allow for creation of a PDF document
** of A4-sized pages, containing a number of images of Doctors, their names,
** positions and some supplementary information
 */
package billedvaeg

import (
	"bytes"
	"sort"
)

// Person contains information on a doctor
type Person struct {
	Name     string
	Position *Position
	Suppl    string
	Img      *bytes.Buffer
}

// SortPersons sorts a slice of persons by positions - and if they're equal,
// by name
func SortPersons(ppl []*Person) []*Person {
	list := make([]*Person, len(ppl))
	copy(list, ppl)
	sort.Slice(list, func(i, j int) bool {
		// Same positions, sort by name instead
		if list[i].Position.Value == list[j].Position.Value {
			return list[i].Name < list[j].Name
		}

		// sort by name
		return list[i].Position.Value < list[j].Position.Value
	})
	return list
}

// Position contains information on a position
type Position struct {
	Title string
	Abbr  string
	Value int
}

/*// PositionsToMap converts a slice of positions into a map,
// where the key is the position's 'abbr' (abbreviation).
func PositionsToMap(pl []*Position) map[string]*Position {
	m := make(map[string]*Position)
	for _, p := range pl {
		m[p.Abbr] = p
	}
	return m
}*/

/*// parseMultiformData parses a HTTP MultiFormData request, creating a PersonList
// from it's information and returning the data.
func parseMultiformData(r *http.Request) (*[]*Person, error) {
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
}*/
