package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialise a new router
	router := httprouter.New()

	// Override defaults error handling on router by replacing them with helper function that satisfy the http.Handler
	// interface.
	router.NotFound = http.HandlerFunc(app.resourceNotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	//================================== MOVIES ======================================================
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMoviesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)

	//================================== USERS =======================================================
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	// ================================ Authentication ===============================================
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(router))
}
