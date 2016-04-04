package thunder

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        )

type Client interface {
  Get(url string) (resp *http.Response, err error)
  Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error)
  Put(url string, bodyType string, body io.Reader) (resp *http.Response, err error)
  Delete(resp *http.Response, err error)
}


type Arrester struct {
  url string
}

func (a Arrester) Set(c Client, key string , value interface, ex) {

}

func (a Arrester) Get(c Client, key string) (interface {}, bool) {

}

func (a Arrester) Update(c Client, key) {

}

func (a Arrester) Del(c Client, key) {

}

func (a Arrester) Keys() {

}
