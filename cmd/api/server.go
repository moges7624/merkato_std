package main

import (
	"fmt"
	"net/http"
)

func (s *APIServer) serve() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: &s.routes,
	}

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
