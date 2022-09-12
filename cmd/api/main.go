package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// version is a string containing the application version number.
const version = "1.0.0"

// config is used to manage the configuration settings of the application.
type config struct {
	port int
	env  string
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
	// Parses the flag values and sets them to the config fields.
	flag.Parse()

	/* Initialize a new logger to write messages to the standard out stream. Logger messages with be prefixed with the
	   time and date.
	*/
	logger := log.New(os.Stdout, "", log.LstdFlags)

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
	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatalf(err.Error())
	}
}
