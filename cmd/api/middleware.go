package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Creates a deferred IIFE to be run in the instance that the server
		// encounters a panic. The panic is dealt with gracefully by sending a
		// 500 response to the client and logging the fatal error encountered.
		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	// Define a client struct to hold the rate limiter and last seen time for each client.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	// Declare variables to hold a map of client IP addresses and associated rate limiters
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// initiate a background goroutine which removes old entries from the clients map, once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			// Loop through all the clients and delete any that haven't been seen in the last 3 minutes
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract client's IP address from the request.
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// lock the mutex to prevent the code being run concurrently.
		mu.Lock()

		// If we don't have the ip address already, create a new limiter for it.
		_, found := clients[ip]
		if !found {
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}
 
		clients[ip].lastSeen = time.Now()

		// If the number of tokens in the limiter bucket is empty, unlock the mutex and return an error
		if !clients[ip].Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}

		// unlock the mutex. We don't defer this as the mutex would only then be unlocked once all handlers downstream
		// of this middleware have also returned.
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
