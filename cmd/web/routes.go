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
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))


	// DYNAMIC MIMIDDLEWARES CHAIN (UNPROTECTED)
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	// STATIC HANDLERS
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))

	// AUTHENTICATION HANDLERS
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))


	// DYNAMIC MIMIDDLEWARES CHAIN (PROTECTED)
	protected := dynamic.Append(app.requireAuthentication)

	// STATIC HANDLERS
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))

	// AUTHENTICATION HANDLERS
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))


	// STANDART MIDDLEWARES CHAIN
	midwares := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return midwares.Then(router)
}