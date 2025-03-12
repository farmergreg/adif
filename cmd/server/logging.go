package main

import (
	"log"
	"net/http"
	"time"
)

// RequestLogger is a middleware that logs information about each HTTP request
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		ipAddress := r.RemoteAddr
		contentType := r.Header.Get("Content-Type")
		bytes := r.Header.Get("Content-Length")

		next.ServeHTTP(w, r)
		log.Printf("%s: %s %s %s %s %s", ipAddress, r.Method, r.URL.Path, contentType, bytes, time.Since(startTime))
	})
}
