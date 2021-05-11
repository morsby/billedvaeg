package billedvaeg

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Serve(port int) {
	r := mux.NewRouter()
	r.HandleFunc("/", UploadHandler).Methods("POST")
	r.HandleFunc("/", Get).Methods("GET")
	http.Handle("/", r)
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	srv := &http.Server{
		Handler: r,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("Listening on http://%s 🚀!\n", addr)
	log.Fatal(srv.ListenAndServe())
}

func Get(w http.ResponseWriter, r *http.Request) {
	Compile(w)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	specialists := r.FormValue("special") == "on"
	fmt.Println(specialists)
	ppl, err := HandleMultiformData(r, FormOptions{MAX_UPLOAD_SIZE: int64(5 << 20), Specialists: specialists})
	if err != nil {
		http.Error(w, fmt.Sprintf("There was an err: %s", err.Error()), http.StatusInternalServerError)
	}
	ppl.Sort()

	doc := New()
	AddPeople(doc, *ppl, 3)

	if err != nil {
		panic(err)
	}
	_ = doc.Output(w)
}
