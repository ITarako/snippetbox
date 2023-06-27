package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.Handler(http.MethodGet, "/", app.sessionManager.LoadAndSave(http.HandlerFunc(app.home)))
	router.Handler(http.MethodGet, "/snippet/view/:id", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetView)))
	router.Handler(http.MethodGet, "/snippet/create", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetCreate)))
	router.Handler(http.MethodPost, "/snippet/create", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetCreatePost)))

	router.Handler(http.MethodGet, "/user/signup", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userSignup)))
	router.Handler(http.MethodPost, "/user/signup", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userSignupPost)))
	router.Handler(http.MethodGet, "/user/login", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLogin)))
	router.Handler(http.MethodPost, "/user/login", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLoginPost)))
	router.Handler(http.MethodPost, "/user/logout", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLogoutPost)))

	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
