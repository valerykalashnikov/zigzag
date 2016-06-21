package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/valerykalashnikov/zigzag/zigzag"
)

type Clock struct {
	ex int64
}

func (c *Clock) Now() time.Time          { return time.Now() }
func (c *Clock) Duration() time.Duration { return time.Duration(c.ex) * time.Minute }

func Index(db *zigzag.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	fmt.Fprintln(w, "ZigZag server!")
	return http.StatusOK, nil
}

func Set(db *zigzag.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	var value interface{}
	var ex int64

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if stat := db.CheckRole(); stat == "slave" {
		return 403, nil
	}

	key := getKey(r)

	ex, err := getExpiration(r.URL.Query())
	if err != nil {
		w.WriteHeader(422) // unprocessable entity
		fmt.Fprintf(w, "Error: Invalid expiration value")
		return 422, err
	}

	body, err := requestBody(r)
	if err != nil {
		return 500, err
	}

	if err := json.Unmarshal(body, &value); err != nil {
		respondWithParsingJsonError(w, err)
		return 422, err
	}

	db.Set(key, value, &Clock{ex})

	w.WriteHeader(http.StatusOK)
	return http.StatusOK, nil
}

func Get(db *zigzag.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	key := getKey(r)

	if item, found := db.Get(key); found {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(item); err != nil {
			panic(err)
		}
		return http.StatusOK, nil
	}
	w.WriteHeader(http.StatusNotFound)
	return http.StatusNotFound, errors.New("Not found")
}

func Update(db *zigzag.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	var value interface{}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if stat := db.CheckRole(); stat == "slave" {
		return 403, nil
	}

	key := getKey(r)

	body, err := requestBody(r)
	if err != nil {
		return 500, err
	}

	if err := json.Unmarshal(body, &value); err != nil {
		respondWithParsingJsonError(w, err)
		return 422, err
	}

	if ok := db.Upd(key, value); ok {
		w.WriteHeader(http.StatusOK)
		return http.StatusOK, nil
	}

	w.WriteHeader(http.StatusNotFound)
	return http.StatusNotFound, errors.New("Not found")

}

func Delete(db *zigzag.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	if stat := db.CheckRole(); stat == "slave" {
		return 403, nil
	}

	key := getKey(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	db.Del(key)
	w.WriteHeader(http.StatusOK)
	return http.StatusOK, nil
}

func Keys(db *zigzag.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	escapedPattern := vars["pattern"]
	pattern, err := url.QueryUnescape(escapedPattern)
	if err != nil {
		w.WriteHeader(422) // unprocessable entity
		fmt.Fprintf(w, "Error: Invalid pattern")
		return 422, err
	}
	keys := db.Keys(pattern)
	if err := json.NewEncoder(w).Encode(keys); err != nil {
		panic(err)
	}
	return http.StatusOK, nil
}

func getExpiration(query url.Values) (int64, error) {
	exStr := query.Get("ex")

	if exStr == "" {
		return 0, nil
	}

	return strconv.ParseInt(exStr, 10, 64)
}

func getKey(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["key"]
}

func requestBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	return body, err
}

func respondWithParsingJsonError(w http.ResponseWriter, err error) {
	w.WriteHeader(422) // unprocessable entity
	fmt.Fprintf(w, "Error: %s", err)
}
