package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

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

	// Execute() is used to write the template to the response body,
	// which the reader receives and renders the template to their browser.
	// The last parameter to Execute represents dynamic data to be rendered
	// in the template
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.serverError(w, r, err)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with id %d", id)
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
