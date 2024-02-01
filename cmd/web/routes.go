package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/julienschmidt/httprouter"
)


func (app *application) routes() http.Handler {
	// ROUTER
	router := httprouter.New()

	// CUSTOM NOT FOUND
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// STATIC
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	router.Handler(http.MethodGet, "/static", http.NotFoundHandler()) // Custom
	router.Handler(http.MethodGet, "/static/", http.StripPrefix("/static", fileServer))

	// HANDLERS
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// MIDDLEWARES
	midwares := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return midwares.Then(router)
}