package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       time.Minute,
		ErrorLog:          log.New(app.logger, "", 0),
	}

	go func() {
		// Create a quit channel that receives os signals
		quit := make(chan os.Signal, 1)

		// Begin listening for any signals on the quit channel that match the interruption or terminate signals.
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// If the quit channel receives a signal, print the signal caught and exit the application with a success
		// status code
		s := <-quit
		app.logger.PrintInfo("caught signal", map[string]string{
			"signal": s.String(),
		})
		// Exit with success status code.
		os.Exit(0)
	}()

	// Start the server.
	app.logger.PrintInfo("server starting", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	return srv.ListenAndServe()
}
