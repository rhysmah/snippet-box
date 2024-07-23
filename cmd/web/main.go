package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Creates a file server that serves files out of the "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home) // Requires exact match
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Printf("starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
