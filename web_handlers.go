package billedvaeg

import (
	"fmt"
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {
	Compile(w)
}

func PostMultiformData(w http.ResponseWriter, r *http.Request) {
	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ppl, err := parseMultiformData(r)
	ppl.Sort()

	if err != nil {
		http.Error(w, fmt.Sprintf("There was an err: %s", err.Error()), http.StatusInternalServerError)
	}
	doc := New()
	AddPeople(doc, ppl, 3)
	_ = doc.Output(w)
}
