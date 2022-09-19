package main

import (
	"fmt"
	"net/http"
	"time"

	"richwynmorris.co.uk/internal/data"
	"richwynmorris.co.uk/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	v.Check(input.Title != "", "title", "title must be provided")
	v.Check(len(input.Title) <= 500, "title", "title must not be more than 500 bytes long")

	v.Check(input.Year != 0, "year", "year must be provided")
	v.Check(input.Year >= 1888, "year", "year must be greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "year cannot be in the future")

	v.Check(input.Runtime != 0, "runtime", "runtime must be provided")
	v.Check(input.Runtime > 0, "runtime", "runtime must be a positive integer")

	v.Check(input.Genres != nil, "genres", "a genre must be selected")
	v.Check(len(input.Genres) >= 1, "genres", "a minimum of one genre must be selected")
	v.Check(len(input.Genres) <= 5, "genres", "There cannot be more than 5 genres selected")
	v.Check(validator.Unique(input.Genres), "genres", "genres cannot contain duplicates")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the params to and integer and check it is a valid integer.
	id, err := app.readIDParam(r)
	if err != nil {
		app.resourceNotFoundResponse(w, r)
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

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
