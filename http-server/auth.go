package main

import (
          "net/http"
          "strings"
          "fmt"
        )

func Auth(handler http.Handler, zigzagToken string) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    token := requestToken(r)
    if zigzagToken == "" || token == zigzagToken {
      handler.ServeHTTP(w, r)
    } else {
      w.WriteHeader(http.StatusUnauthorized)
      return
    }
  })
}

func requestToken(r *http.Request) string {
  authStr := r.Header.Get("Authorization")
  fmt.Println(authStr)
  if !strings.HasPrefix(authStr, "Token ") {
    return ""
  }
  return authStr[6:]
}
