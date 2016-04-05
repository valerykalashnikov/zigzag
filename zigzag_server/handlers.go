package main

import (
  "time"
  "encoding/json"
  "fmt"
  "net/http"
  "net/url"
  "io"
  "io/ioutil"

  "strconv"

  "github.com/gorilla/mux"
  "zigzag/zigzag"
)

type Clock struct{
  ex int64
}

func (c *Clock) Now() time.Time { return time.Now() }
func (c *Clock) Duration() time.Duration { return time.Duration(c.ex) * time.Minute }

func Index(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "ZigZag server!")
}

func Set(w http.ResponseWriter, r *http.Request) {
  var value interface {}
  var ex int64

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")

  key := getKey(r)

  ex, err := getExpiration(r.URL.Query())
  if err != nil {
    w.WriteHeader(422) // unprocessable entity
    fmt.Fprintf(w, "Error: Invalid expiration value")
    return
  }

  body, err := requestBody(r);
  if err != nil {panic(err)}

  if err := json.Unmarshal(body, &value); err != nil {
    respondWithParsingJsonError(w, err)
  }

  zigzag.Set(key, value, &Clock{ex})

  w.WriteHeader(http.StatusOK)
}

func Get(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")

  key := getKey(r)

  if item, found := zigzag.Get(key); found {
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(item); err != nil {
        panic(err)
    }
    return
  }
  w.WriteHeader(http.StatusNotFound)
}

func Update(w http.ResponseWriter, r *http.Request) {
  var value interface {}

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")

  key  := getKey(r)

  body, err := requestBody(r)
  if err != nil {panic(err)}

  if err := json.Unmarshal(body, &value); err != nil {
    respondWithParsingJsonError(w, err)
    return
  }

  if ok := zigzag.Upd(key, value); ok {
    w.WriteHeader(http.StatusOK)
    return
  }

  w.WriteHeader(http.StatusNotFound)

}

func Delete(w http.ResponseWriter, r *http.Request) {
  key := getKey(r)
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  zigzag.Del(key)
  w.WriteHeader(http.StatusOK)
}

func Keys(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  escapedPattern := vars["pattern"]
  pattern, err := url.QueryUnescape(escapedPattern)
  if (err != nil) {
    w.WriteHeader(422) // unprocessable entity
    fmt.Fprintf(w, "Error: Invalid pattern")
    return
  }
  keys := zigzag.Keys(pattern)
  if err := json.NewEncoder(w).Encode(keys); err != nil {
    panic(err)
  }
}

func getExpiration(query url.Values) (int64, error) {
  exStr := query.Get("ex")

  if exStr == "" { return 0, nil }

  return strconv.ParseInt(exStr, 10, 64)
}

func getKey(r *http.Request) string {
  vars := mux.Vars(r)
  return vars["key"]
}

func requestBody(r *http.Request) ([]byte, error) {
  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

  if err != nil { panic(err) }
  if err := r.Body.Close(); err != nil { panic(err) }
  return body, err
}

func respondWithParsingJsonError(w http.ResponseWriter, err error) {
  w.WriteHeader(422) // unprocessable entity
  fmt.Fprintf(w, "Error: %s", err)
}

