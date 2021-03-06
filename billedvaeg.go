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
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	PositionID    int           `json:"positionId"`
	PositionOrder int           `json:"position"`
	Suppl         string        `json:"suppl"`
	Img           *bytes.Buffer `json:"image"`
	Order         int           `json:"order"`
}

type personJson struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	PositionID    int        `json:"positionId"`
	PositionOrder int        `json:"position"`
	Suppl         string     `json:"suppl"`
	Order         int        `json:"order"`
	Img           *base64Img `json:"image"`
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

	p.ID = tmp.ID
	p.Name = tmp.Name
	p.PositionID = tmp.PositionID
	p.PositionOrder = tmp.PositionOrder
	p.Suppl = tmp.Suppl
	p.Order = tmp.Order
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
func SortPersons(ppl []*Person, auto bool) []*Person {
	list := make([]*Person, len(ppl))
	copy(list, ppl)
	sort.Slice(list, func(i, j int) bool {
		if auto {
			// Same positions, sort by name instead
			if list[i].PositionOrder == list[j].PositionOrder {
				return list[i].Name < list[j].Name
			}

			// sort by position
			return list[i].PositionOrder < list[j].PositionOrder
		} else {
			return list[i].Order < list[j].Order
		}

	})
	return list
}

// Position contains information on a position
type Position struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Abbr  string `json:"abbr"`
	Order int    `json:"order"`
}
