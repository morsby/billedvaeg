package billedvaeg

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Server starts a server listening in PORT.
// wait is the duration for which the server gracefully waits for existing connections to finish - e.g. 15s or 1m
func Server(port int, wait time.Duration) {

	r := mux.NewRouter()
	r.HandleFunc("/", post).Methods("POST")
	r.HandleFunc("/", get).Methods("GET")
	// Add your routes as needed

	// allow CORS
	handler := cors.Default().Handler(r)

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		fmt.Printf("Server listening on http://%s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Please issue a post request"))
}

func post(w http.ResponseWriter, r *http.Request) {

	var input JSONInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		panic(err)
	}

	doc := New()
	doc.People = input.People
	doc.Positions = input.Positions
	err = doc.Generate(input.Sort)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/pdf")
	encoder := base64.NewEncoder(base64.StdEncoding, w)
	err = doc.PDF.OutputAndClose(encoder)

	if err != nil {
		panic(err)
	}
}
