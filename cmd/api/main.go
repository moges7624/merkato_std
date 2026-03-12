package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type APIServer struct {
	port   int
	routes http.ServeMux
	logger *slog.Logger
	dsn    string
}

func main() {
	var port int
	var dsn string

	flag.IntVar(&port, "port", 4000, "API server port")
	flag.StringVar(
		&dsn,
		"db-dsn",
		"postgres://merkato:123456@localhost/merkatostd?sslmode=disable",
		"PostgreSQL DSN",
	)

	logHanlder := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logHanlder)

	db, err := openDB(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &APIServer{
		port:   port,
		logger: logger,
	}

	app.routes = *app.NewRouter()

	app.logger.Info("Starting server", "port", app.port)

	err = app.serve()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
