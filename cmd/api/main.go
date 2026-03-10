package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type APIServer struct {
	port   int
	routes http.ServeMux
	logger *slog.Logger
}

func main() {
	var port int

	flag.IntVar(&port, "port", 4000, "API server port")

	logHanlder := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logHanlder)

	app := &APIServer{
		port:   port,
		logger: logger,
	}

	app.routes = *app.NewRouter()

	app.logger.Info("Starting server", "port", app.port)

	err := app.serve()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}
