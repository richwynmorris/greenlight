package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Print(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"message": message}

	err := app.writeJSON(w, status, env, nil)
	if nil != err {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	msg := "The server is unable to fulfil this request at this time. Please try again later."
	app.errorResponse(w, r, http.StatusInternalServerError, msg)
}

func (app *application) resourceNotFoundResponse(w http.ResponseWriter, r *http.Request) {
	msg := "The resource requested could not be found"
	app.errorResponse(w, r, http.StatusNotFound, msg)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("The method %s is not allowed for this request", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, msg)
}
