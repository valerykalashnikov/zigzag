package main

import (
    "log"
    "net/http"
    "zigzag/jobs"
)

func main() {

    go jobs.CleanCache()

    router := NewRouter()

    log.Fatal(http.ListenAndServe(":8080", router))
}


