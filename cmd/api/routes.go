package main

import (
	"io"
	"net/http"
)

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "hello there")
}

func (s *APIServer) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", homeHandler)

	return mux
}
