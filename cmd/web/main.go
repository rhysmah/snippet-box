package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/rhysmah/snippet-box/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	// We're not using anything from this import, so we prefix it with an
	// underscore, else we'll get a compile-time error. We need the `init`
	// function to run from this package so it can register itself with the
	// database/sql package.
	_ "github.com/go-sql-driver/mysql"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now, include only the structured logger; more to be added.
// Add the SnippetModel from the `internal/models` directory; like the logger,
// we've injected this as a dependency in our application.
type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

// TODO (if applicable): create a `config` struct for configuration settings

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:1234@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Database
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	// Initialize new session manager; configure it to
	// use MySQL database as the session store. Set a
	// lifetime of 12 hours, meaning sessions will expire
	// 12 hours after the sessions are created and stored.
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	server := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Because `addr` is a pointer to a string, we must
	// dereference it to get access to its value.
	// The first argument is the message; the variadic variables that
	// follow are key-value pairs.
	// Using slog.String() is optional, but it adds a level of type-safety
	// and prevents, for example, leaving out a key or value.
	logger.Info("starting server", slog.String("addr", *addr))

	err = server.ListenAndServe()

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
