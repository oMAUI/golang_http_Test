package main

import (
	"errors"
	"net/http"
	"refactoring/HttpRoute"
)

const store = `users.json`

var (
	UserNotFound = errors.New("user_not_found")
)

func main() {
	route := HttpRoute.NewRoute()
	http.ListenAndServe(":3333", route)
}