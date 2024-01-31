package main

import "net/http"


func (app *application) routes() *http.ServeMux {
	// SERVEMUX
	mux := http.NewServeMux()

	// STATIC
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux.Handle("/static", http.NotFoundHandler()) // Custom
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// HANDLERS
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}