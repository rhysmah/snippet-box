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
	// flag.Parse() reads CL flag value and assigns it to `addr` variable.
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Use the slog.New() function to initialize a new structured (custom) logger
	// For now, it'll write to the stdout and use the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // Adds filename and line number to log
	}))

	app := &application{
		logger: logger,
	}

	// Because `addr` is a pointer to a string, we must
	// dereference it to get access to its value.
	// The first argument is the message; the variadic variables that
	// follow are key-value pairs.
	// Using slog.String() is optional, but it adds a level of type-safety
	// and prevents, for example, leaving out a key or value.
	logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
