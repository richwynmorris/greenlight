package main

import (
	"fmt"
	"net/http"
	"time"

	"richwynmorris.co.uk/internal/data"
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

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablance",
		Runtime:   103,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process you request", http.StatusBadRequest)
	}
}
