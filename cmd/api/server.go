package main

import (
	"context"
	"errors"
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

	shutdownErr := make(chan error)

	go func() {
		// Create a quit channel that receives os signals
		quit := make(chan os.Signal, 1)

		// Begin listening for any signals on the quit channel that match the interruption or terminate signals.
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// If the quit channel receives a signal, print the signal caught
		s := <-quit
		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		// create a context to provide 20 seconds to end any connections
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		shutdownErr <- srv.Shutdown(ctx)
	}()

	// Start the server.
	app.logger.PrintInfo("server starting", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	// Check if the graceful shutdown process had begun. If it has, and the error is not the ErrServerClosed,
	// we need to return the error
	err := srv.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	// wait to receive an error from the shutdown process. If there was an error in gracefully shutting down
	// then we need to return the error:
	err = <-shutdownErr
	if err != nil {
		return err
	}

	// At this point we know that the graceful shutdown completed successfully.
	app.logger.PrintInfo("stopped server", map[string]string{"addr": srv.Addr})

	return nil
}
