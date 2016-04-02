package main

import (
  "net/http"
  "github.com/gorilla/mux"
)

func NewRouter(authToken string) *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
        var handler http.Handler
        handler = route.HandlerFunc
        handler = Auth(Logger(handler, route.Name), authToken)

        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)

    }
    return router
}
