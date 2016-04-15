package main

type Route struct {
  Name        string
  Method      string
  Pattern     string
  HandlerFunc AppHandlerFunc
}

type Routes []Route

var routes = Routes{
  Route{
      "index",
      "GET",
      "/",
      Index,
  },
  Route{
      "keys",
      "GET",
      "/keys/{pattern}",
      Keys,
  },
  Route{
      "set",
      "POST",
      "/set/{key}",
      Set,
  },
  Route{
      "get",
      "GET",
      "/get/{key}",
      Get,
  },
  Route{
      "update",
      "PUT",
      "/update/{key}",
      Update,
  },
  Route{
      "delete",
      "DELETE",
      "/delete/{key}",
      Delete,
  },
}
