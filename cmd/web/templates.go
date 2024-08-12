package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/rhysmah/snippet-box/internal/models"
)

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Use `filepath.Glob()` to get a slice of all filepaths
	// that match the pattern. This gives us a slice of all the
	// filepaths for out application 'page' templates.
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the filename (e.g., `home.tmpl.html`) from
		// the full filepath and assign it to `name` variable
		name := filepath.Base(page)

		// Parse the files into a template set
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add template set to map; use name of page
		// (e.g. `home.tmpl.html`) as the key.
		cache[name] = ts
	}

	return cache, nil
}

// Define a templateData type to act as the holding
// structure for any dynamic data that we want to pass
// to our HTML templates. We can only pass one piece
// of dynamic data, so a struct is a way to contain
// one datum composed of many data.
type templateData struct {
	Year            int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}
