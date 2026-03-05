package main

import (
	"flag"
	"fmt"
	"net/http"
)

type APIServer struct {
	port   int
	routes http.ServeMux
}

func main() {
	var port int

	flag.IntVar(&port, "port", 4000, "API server port")

	app := &APIServer{
		port: port,
	}

	app.routes = *app.NewRouter()

	fmt.Printf("Starting server on port %d...", app.port)

	err := app.serve()
	if err != nil {
		panic(err)
	}
}
