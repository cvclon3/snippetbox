package main

import (
	"net/http"

	"snippetbox.cvclon3.net/ui"

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
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// PING
	router.HandlerFunc(http.MethodGet, "/ping", ping)


	// DYNAMIC MIMIDDLEWARES CHAIN (UNPROTECTED)
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// STATIC HANDLERS
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.about))

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
	router.Handler(http.MethodGet, "/account/view", protected.ThenFunc(app.accountView))


	// STANDART MIDDLEWARES CHAIN
	midwares := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return midwares.Then(router)
}