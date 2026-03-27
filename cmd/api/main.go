package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/moges7624/merkato_std/internal/config"
)

type APIServer struct {
	routes http.ServeMux
	logger *slog.Logger
	cfg    config.Config
	DB     *sql.DB
}

func main() {
	logHanlder := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logHanlder)

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("error loading config", "detail", err.Error())
		os.Exit(1)
	}

	db, err := openDB(cfg.DSN)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	if err != nil {
		logger.Error("invalid port number")
		os.Exit(1)
	}

	app := &APIServer{
		logger: logger,
		cfg:    *cfg,
		DB:     db,
	}

	app.routes = *app.NewRouter()

	app.logger.Info("Starting server", "port", app.cfg.Port)

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
		db.Close()
		return nil, err
	}

	return db, nil
}
