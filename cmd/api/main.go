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

	"richwynmorris.co.uk/internal/data"
	"richwynmorris.co.uk/internal/jsonlog"
)

// version is a string containing the application version number.
const version = "1.0.0"

// config is used to manage the configuration settings of the application.
type config struct {
	port int
	env  string
	db   struct {
		dsn         string
		maxOpenConn int
		maxIdleConn int
		maxIdleTime string
	}
}

// application holds the handlers, helpers and middleware to support the application's functionality.
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config

	// Initialise flag names and default values to run the application.
	flag.IntVar(&cfg.port, "port", 4000, "API Server Port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")

	// Use flags to set database connection pool settings.
	flag.IntVar(&cfg.db.maxOpenConn, "db-max-open-connections", 25, "Max open connections for database")
	flag.IntVar(&cfg.db.maxIdleConn, "db-max-idle-connections", 25, "Max idle conections for database")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "Max idle time for database")

	// Parses the flag values and sets them to the config fields.
	flag.Parse()

	/* Initialize a new logger to write messages to the standard out stream. Logger messages with be prefixed with the
	   time and date.
	*/
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	// Declare instance of application struct with logger and config settings.
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
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
		ErrorLog:          log.New(logger, "", 0),
	}

	// Start the server.
	logger.PrintInfo("server starting", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.env,
	})
	err = srv.ListenAndServe()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	// Create an empty connection pool using the dsn config.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	parsedTime, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConn)
	db.SetConnMaxIdleTime(parsedTime)
	db.SetMaxIdleConns(cfg.db.maxIdleConn)

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
