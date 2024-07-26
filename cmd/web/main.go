package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	// We're not using anything from this import, so we prefix it with an
	// underscore, else we'll get a compile-time error. We need the `init`
	// function to run from this package so it can register itself with the
	// database/sql package.
	_ "github.com/go-sql-driver/mysql"
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

	// CL flag for the MySQL DSN string.
	dsn := flag.String("dsn", "web:1234@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Use the slog.New() function to initialize a new structured (custom) logger
	// For now, it'll write to the stdout and use the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // Adds filename and line number to log
	}))

	// Database
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

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

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {

	// This does NOT open any connections; it simply initializes
	// a pool of connections for future use. These connections are
	// Go manages this pool of connections automatically, opening
	// and closing connections to the database via the driver.
	// This pool of connections is safe for concurrent use.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Connections to the DB are established lazily -- as ands
	// when needed. Thus, to check if the database is accessible,
	// we need to ping it.
	err = db.Ping()
	if err != nil {
		db.Close() // If an error, immediately close the DB.
		return nil, err
	}

	return db, nil
}
