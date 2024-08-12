package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
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

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		Year:            time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {

	// Call ParseForm on the request
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Attempt to decode our form; pass the target destination, `dst`, as first param
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {

		// If we attempt to use an invalid target destination, Decode() will
		// return an error with type *form.InvalidDecoderError
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
