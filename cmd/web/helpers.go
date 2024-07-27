package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError() helper writes log entry at Error level (including request method and URI as atts),
// then sends generic 500 Internal Service Error response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError() helper sends specific status code and corresponding description to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {

	// Retrieve the appropriate template set from the cached based on page name (e.g., `home.tmpl.html`).
	// If no entry exists in cache, create a new error and call serverError() helper method
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s doesn't exist", page)
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	// Execute template set and write the response body;
	// if any error, call serverError() helper.
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
