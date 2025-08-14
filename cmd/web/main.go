package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"snippetbox.siddharthroy.com/internal/models"
)

type config struct {
	addr      string
	staticDir string
}

type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecorder   *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	loggerHandler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(loggerHandler)

	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	db, err := openDB(*dsn)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecorder := form.NewDecoder()

	sessionManager := scs.New()

	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := application{
		logger: logger,
		snippets: &models.SnippetModel{
			DB: db,
		},
		users: &models.UserModel{
			DB: db,
		},
		templateCache:  templateCache,
		formDecorder:   formDecorder,
		sessionManager: sessionManager,
	}

	logger.Info("starting server", "addr", cfg.addr)

	srv := &http.Server{
		Addr:     cfg.addr,
		Handler:  app.routes(&cfg),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	err = srv.ListenAndServe()

	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

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
