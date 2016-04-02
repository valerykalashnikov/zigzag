package main

import (
  "fmt"
  "net"
  "encoding/gob"
  "zigzag/zigzag"
)

type P struct {
  M, N int64
}
func handleConnection(conn net.Conn) {
  dec := gob.NewDecoder(conn)
  p := &P{}
  dec.Decode(p)
  fmt.Printf("Received : %+v", p);
}

func main() {
  fmt.Println("ZigZag server started:");
  ln, err := net.Listen("tcp", ":8081")
  if err != nil {
    fmt.Printf("Error : %+v", err);
  }
  for {
      conn, err := ln.Accept() // this blocks until connection or error
      if err != nil {
          fmt.Printf("Error : %+v", err);
          continue
      }
      go handleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
  }
}
