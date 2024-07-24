package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now, include only the structured logger; more to be added.
type application struct {
	logger *slog.Logger
}

// TODO (if applicable): create a `config` struct for configuration settings

func main() {

	// Define a new command-line flag so we can set the network address.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// This reads the the command-line flag value and assigns it to the `addr` variable
	flag.Parse()

	// Use the slog.New() function to initialize a new structured (custom) logger
	// For now, it'll write to the stdout and use the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // Adds filename and line number to log
	}))

	app := &application{
		logger: logger,
	}

	mux := http.NewServeMux()

	// Creates a file server that serves files out of the "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.home) // Requires exact match
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Because `addr` is a pointer to a string, we must
	// dereference it to get access to its value.
	// The first argument is the message; the variadic variables that
	// follow are key-value pairs.
	// Using slog.String() is optional, but it adds a level of type-safety
	// and prevents, for example, leaving out a key or value.
	logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
