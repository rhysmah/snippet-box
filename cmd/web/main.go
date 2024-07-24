package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	// Define a new command-line flag so we can set the network address.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// This reads the the command-line flag value and assigns it to the `addr` variable
	flag.Parse()

	mux := http.NewServeMux()

	// Creates a file server that serves files out of the "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home) // Requires exact match
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// Because `addr` is a pointer to a string, we must
	// dereference it to get access to its value.
	log.Printf("starting server on %s", *addr)

	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
