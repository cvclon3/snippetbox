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

	// DYNAMIC MIMIDDLEWARES CHAIN
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// HANDLERS
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// STANDART MIDDLEWARES CHAIN
	midwares := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return midwares.Then(router)
}