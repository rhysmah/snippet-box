package main

import (
	"bytes"
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

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {

	// Retrieve the appropriate template set from the cached based on page name (e.g., `home.tmpl.html`).
	// If no entry exists in cache, create a new error and call serverError() helper method
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s doesn't exist", page)
		app.serverError(w, r, err)
		return
	}

	// We'll write the page to a buffer; if there's no issue,
	// then we can write the data to the response the user
	// receives, else we send them an error.
	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	// If there's no error, write data from buffer to response
	buf.WriteTo(w)
}
