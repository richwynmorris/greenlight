package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

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

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {
	query := `
	INSERT INTO movies (title, year, runtime, genres)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, created_at, title, year, runtime, genres, version FROM movies
			  WHERE id = $1`

	var movie Movie

	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}
func (m MovieModel) Update(movie *Movie) error {
	query := `UPDATE movies
			  SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
			  WHERE id = $5
              RETURNING version`

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	// QueryRow expects to return a single row from the db, if it doesn't it throws and error.
	// Query row takes two args: the query and the values to be interpolated into the query.
	// We do this to prevent SQL injection attacks.
	// We also use the rest syntax to explode the values in the slice.
	// The Scan method receives the result of the query and copies the return values into the destination argument.
	err := m.DB.QueryRow(query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	// Return nil as success if update operation performed correctly.
	return nil
}
func (m MovieModel) Delete(id int64) error {
	return nil
}
