package main

import (
    "fmt"
    "log"
    "net"
    "encoding/gob"
    "zigzag/zigzag"
)

type Clock struct{
  ex int64
}

func (c *Clock) Now() time.Time { return time.Now() }
func (c *Clock) Duration() time.Duration { return time.Duration(c.ex) * time.Minute }

var conn string

func Set(key, value, ex) {
  encoder := gob.NewEncoder(conn)
}

func init(){
  conn, err := net.Dial("tcp", "localhost:8080")
  if err != nil {
      log.Fatal("Connection error", err)
  }
}

func main() {
  fmt.Println("start client");
  conn, err := net.Dial("tcp", "localhost:8081")
  if err != nil {
      log.Fatal("Connection error", err)
  }
  encoder := gob.NewEncoder(conn)
  cache.Set()
  encoder.Encode(p)
  conn.Close()
  fmt.Println("done");
}
