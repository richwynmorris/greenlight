package main

import (
	"fmt"
	"net/http"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Creating a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the params to and integer and check it is a valid integer.
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Showing Movie with ID: %i\n", id)
}
