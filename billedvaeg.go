/* Package billedvaeg contains an api to allow for creation of a PDF document
** of A4-sized pages, containing a number of images of Doctors, their names,
** positions and some supplementary information
 */
package billedvaeg

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"sort"
)

type JSONInput struct {
	Positions []*Position `json:"positions"`
	People    []*Person   `json:"people"`
	Sort      bool        `json:"sort"`
}

// Person contains information on a doctor
type Person struct {
	Name     string        `json:"name"`
	Position int           `json:"position"`
	Suppl    string        `json:"suppl"`
	Img      *bytes.Buffer `json:"img"`
}

type personJson struct {
	Name     string     `json:"name"`
	Position int        `json:"position"`
	Suppl    string     `json:"suppl"`
	Img      *base64Img `json:"img"`
}

type base64Img struct {
	MIME string `json:"mime"`
	Data string `json:"data"`
}

func (p *Person) UnmarshalJSON(data []byte) error {
	tmp := personJson{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	p.Name = tmp.Name
	p.Position = tmp.Position
	p.Suppl = tmp.Suppl
	if tmp.Img != nil {
		err := p.ImageFromBase64(tmp.Img.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

// ImageFromReader sets the person's image from the provided io.Reader.
func (p *Person) ImageFromReader(r io.Reader) error {
	buf := bytes.Buffer{}
	buf.ReadFrom(r)
	p.Img = &buf
	return nil
}

// ImageFromFile sets the person's image by loading the image at the provided path.
func (p *Person) ImageFromFile(path string) error {
	img, err := os.Open(path)
	if err != nil {
		return err
	}
	defer img.Close()
	p.ImageFromReader(img)
	return nil
}

// ImageFromBase64 sets the person's image by parsing a base64-encoded image passed as
// a string parameter.
func (p *Person) ImageFromBase64(ImgBase64 string) error {
	// decode base64 encoded image into a buffer
	img, err := base64.StdEncoding.DecodeString(ImgBase64)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte(img))
	p.Img = buf
	return nil
}

// SortPersons sorts a slice of persons by positions - and if they're equal,
// by name
func SortPersons(ppl []*Person) []*Person {
	list := make([]*Person, len(ppl))
	copy(list, ppl)
	sort.Slice(list, func(i, j int) bool {
		// Same positions, sort by name instead
		if list[i].Position == list[j].Position {
			return list[i].Name < list[j].Name
		}

		// sort by position
		return list[i].Position < list[j].Position
	})
	return list
}

// Position contains information on a position
type Position struct {
	Title string `json:"title"`
	Abbr  string `json:"abbr"`
	Value int    `json:"value"`
}
