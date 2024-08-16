package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/rhysmah/snippet-box/internal/models"
	"github.com/rhysmah/snippet-box/ui"
)

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

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
