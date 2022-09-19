package data

import (
	"time"

	"richwynmorris.co.uk/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id,omitempty"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, input *Movie) {
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
}
