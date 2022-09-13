package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Test string with interpolated values using the %q verb.
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.logger.Printf(err.Error())
		http.Error(w, "failed to marshall json data and formulate response", http.StatusBadRequest)
		return
	}
}
