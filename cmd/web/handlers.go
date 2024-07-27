package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/rhysmah/snippet-box/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	webFiles := []string{
		"./ui/html/base.tmpl.html", // base template must be first
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	// The template.ParseFiles() function reads the template and creates
	// a template.Template object, which has template-specific functions that
	// allow us to interact and execute the template. If there's an error
	// parsing the template, we log the error, then send the user an error
	// message with the "Internal Server Error" error code.
	ts, err := template.ParseFiles(webFiles...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := templateData{
		Snippets: snippets,
	}

	// Execute() is used to write the template to the response body,
	// which the reader receives and renders the template to their browser.
	// The last parameter to Execute represents dynamic data to be rendered
	// in the template
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Use SnippetModel's Get() method to retrieve data for specific record based on ID.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Initialized slice containing paths to view.tmpl.html file.
	// This includes the base layout and navication partial.
	webFiles := []string{
		"./ui/html/base.tmpl.html", // base template must be first
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/view.tmpl.html",
	}

	// The template.ParseFiles() function reads the template and creates
	// a template.Template object, which has template-specific functions that
	// allow us to interact and execute the template. If there's an error
	// parsing the template, we log the error, then send the user an error
	// message with the "Internal Server Error" error code.
	ts, err := template.ParseFiles(webFiles...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := templateData{
		Snippet: snippet,
	}

	// Execute() is used to write the template to the response body,
	// which the reader receives and renders the template to their browser.
	// The last parameter to Execute represents dynamic data to be rendered
	// in the template
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "0 snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
