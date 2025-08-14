package main

import (
	"net/http"

	"snippetbox.siddharthroy.com/ui"
)

func (app *application) routes(conf *config) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServerFS(ui.Files)

	mux.Handle("GET /static/", fileServer)

	mux.HandleFunc("GET /{$}", app.HomePage)
	mux.HandleFunc("GET /", app.NotFoundPage)
	mux.HandleFunc("GET /snippet/view/{id}", app.ViewSnippetPage)
	mux.Handle("GET /snippet/create", app.requireAuthentication(http.HandlerFunc(app.CreateSnippetPage)))
	mux.Handle("POST /snippet/create", app.requireAuthentication(http.HandlerFunc(app.CreateSnippet)))
	// Auth pages
	mux.HandleFunc("GET /user/signup", app.SignupPage)
	mux.HandleFunc("POST /user/signup", app.Signup)
	mux.HandleFunc("GET /user/login", app.LoginPage)
	mux.HandleFunc("POST /user/login", app.Login)
	mux.Handle("POST /user/logout", app.requireAuthentication(http.HandlerFunc(app.Logout)))

	return app.recoverPanic(app.logRequests(app.commonHeader(app.authenticate(app.saveAndLoadSession(mux)))))
}
