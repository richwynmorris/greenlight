package main

import (
	"expvar"
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

	//=================================== HEALTH & METRICS ====================================================

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	//================================== MOVIES ======================================================

	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermissions("movies:write", app.deleteMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermissions("movies:read", app.listMoviesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermissions("movies:read", app.showMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermissions("movies:write", app.updateMovieHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermissions("movies:write", app.createMovieHandler))

	//================================== USERS =======================================================

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	// ================================ AUTHENTICATION ===============================================

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// =============================== MIDDLEWARE ===================================================

	// Panic Recovery; Enable Cors; Rate Limiting; Authentication.
	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authentication(router)))))
}
