package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form/v4"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("this template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}
	var buf bytes.Buffer

	err := ts.ExecuteTemplate(&buf, "base", data)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
}

func (app *application) show404(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/404", http.StatusSeeOther)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecorder.Decode(dst, r.PostForm)

	if err != nil {
		var invalidDecodeError *form.InvalidEncodeError

		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}

		return err
	}
	return nil
}
