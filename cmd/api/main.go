package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// version is a string containing the application version number.
const version = "1.0.0"

// config is used to manage the configuration settings of the application.
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

// application holds the handlers, helpers and middleware to support the application's functionality.
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	// Initialise flag names and default values to run the application.
	flag.IntVar(&cfg.port, "port", 4000, "API Server Port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")

	// Parses the flag values and sets them to the config fields.
	flag.Parse()

	/* Initialize a new logger to write messages to the standard out stream. Logger messages with be prefixed with the
	   time and date.
	*/
	logger := log.New(os.Stdout, "", log.LstdFlags)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	defer db.Close()

	logger.Printf("database connection pool established")

	// Declare instance of application struct with logger and config settings.
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Declare a new http-router to dispatch requests to tha application's handlers.
	router := app.routes()

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.port),
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       time.Minute,
	}

	// Start the server.
	logger.Printf("Starting %s server on port: %d", cfg.env, cfg.port)
	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatalf(err.Error())
	}
}

func openDB(cfg config) (*sql.DB, error) {
	// Create an empty connection pool using the dsn config.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Create a context that will time out after 5 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	// Pass the context to the PingContext function to check that the db is working correctly
	// and can be connected to within 5 seconds.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
