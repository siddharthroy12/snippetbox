package main

import "net/http"

func (app *application) routes(conf *config) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(conf.staticDir))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.HomePage)
	mux.HandleFunc("GET /", app.NotFoundPage)
	mux.HandleFunc("GET /snippet/view/{id}", app.ViewSnippetPage)
	mux.HandleFunc("GET /snippet/create", app.CreateSnippetPage)
	mux.HandleFunc("POST /snippet/create", app.CreateSnippet)

	return app.recoverPanic(app.logRequests(app.commonHeader(mux)))
}
