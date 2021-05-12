package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/morsby/billedvaeg"
)

func main() {
	port := 3000
	r := mux.NewRouter()
	r.HandleFunc("/", billedvaeg.PostMultiformData).Methods("POST")
	r.HandleFunc("/", billedvaeg.Get).Methods("GET")
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
