package main

import (
  "log"
  "net/http"
  "github.com/valerykalashnikov/zigzag/zigzag"
  "github.com/gorilla/mux"
)

type AppHandlerFunc func(*zigzag.DB, http.ResponseWriter, *http.Request) (int, error)

type AppHandler struct {
  db *zigzag.DB
  H AppHandlerFunc
}

func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  status, err := ah.H(ah.db, w, r)

  if err != nil { log.Printf("HTTP %d: %q", status, err) }
}

func NewRouter(authToken string, db *zigzag.DB) *mux.Router {
  router := mux.NewRouter().StrictSlash(true)
  for _, route := range routes {
    var handler http.Handler
    handler = AppHandler{db, route.HandlerFunc}
    handler = Auth(Logger(handler, route.Name), authToken)

    router.
      Methods(route.Method).
      Path(route.Pattern).
      Name(route.Name).
      Handler(handler)

  }
  return router
}