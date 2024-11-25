package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	port := "8080"

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	log.Fatal(srv.ListenAndServe())
}
