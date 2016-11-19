package main

import (
	"log"
	"net/http"
	"strings"
)

func Auth(handler http.Handler, zigzagToken string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := requestToken(r)
		if zigzagToken == "" || token == zigzagToken {
			handler.ServeHTTP(w, r)
		} else {
			handleUnauthorized(w)
			return
		}
	})
}

func handleUnauthorized(w http.ResponseWriter) {
	unauthorizedStatus := http.StatusUnauthorized
	w.WriteHeader(unauthorizedStatus)
	log.Printf("HTTP %d: Unathorized", unauthorizedStatus)
}

func requestToken(r *http.Request) string {
	authStr := r.Header.Get("Authorization")
	if !strings.HasPrefix(authStr, "Token ") {
		return ""
	}
	return authStr[6:]
}
