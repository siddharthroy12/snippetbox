package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.siddharthroy.com/internal/models"
	"snippetbox.siddharthroy.com/internal/validator"
)

func (app *application) HomePage(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, r, 200, "home.html", data)

}

func (app *application) ViewSnippetPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.show404(w, r)
			return
		}
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, 200, "view.html", data)

}

func (app *application) CreateSnippetPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = CreateSnippetForm{
		Expires: 365,
	}
	app.render(w, r, http.StatusOK, "create.html", data)
}

type CreateSnippetForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) CreateSnippet(w http.ResponseWriter, r *http.Request) {

	var form CreateSnippetForm

	err := app.decodePostForm(r, &form)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValues(form.Expires, 1, 3, 7), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) NotFoundPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Write([]byte("You are lost"))
}
