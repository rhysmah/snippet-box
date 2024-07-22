package main

import (
	"log"
	"net/http"
)

// home() is a regular Go function with two parameters
//
// Parameters:
//
//	w http.ResponseWriter: provides methods for assembling an HTTP response
//		and sending it to the user.
//	r *http.Request: pointer to a struct which holds information about the
//		current request (like the HTTP method and the requested URL)
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func main() {

	// A ServeMux is an HTTP request multiplexer.
	// It's a map of patterns, strings, to handler functions.
	// When an incoming user request matches a registered pattern, it
	// calls and executes the mapped handler function.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Println("Starting server on :4000")

	// The web server will listen on port 4000.
	// It uses the ServerMux -- named `mux`, in this case --
	// to handle incoming patterns. When a request comes in,
	// it'll check them against patterns registered with `mux`.
	// If there's a handler function assoicated with that pattern,
	// that handler function is executed.
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
