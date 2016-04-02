package main

import (
    "log"
    "net/http"
    "zigzag/jobs"
    "os"
)

func main() {
  port, authToken := os.Getenv("ZIGZAG_PORT"), os.Getenv("ZIGZAG_AUTH")
  if port == "" { port = "8082" }

  go jobs.CleanCache(20)

  router := NewRouter(authToken)

  log.Fatal(http.ListenAndServe(":" + port, router))
}


