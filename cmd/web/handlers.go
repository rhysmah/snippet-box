package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rhysmah/snippet-box/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v", snippet)
	}

	// webFiles := []string{
	// 	"./ui/html/base.tmpl.html", // base template must be first
	// 	"./ui/html/partials/nav.tmpl.html",
	// 	"./ui/html/pages/home.tmpl.html",
	// }

	// // The template.ParseFiles() function reads the template and creates
	// // a template.Template object, which has template-specific functions that
	// // allow us to interact and execute the template. If there's an error
	// // parsing the template, we log the error, then send the user an error
	// // message with the "Internal Server Error" error code.
	// ts, err := template.ParseFiles(webFiles...)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// // Execute() is used to write the template to the response body,
	// // which the reader receives and renders the template to their browser.
	// // The last parameter to Execute represents dynamic data to be rendered
	// // in the template
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// }
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Use the SnippetModel's Get() method to retrieve data for a
	// specific records based on its ID. If no matching reords is
	// found, return a 404 Not Found response.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	fmt.Fprintf(w, "%+v", snippet)
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
