package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RexKizzy22/go-movies-api/models"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type DB struct {
	dsn string
}
type JWT struct {
	secret string
}

type config struct {
	port int
	env  string
	db   DB
	jwt  JWT
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// var cfg config

	// flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	// flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production)")
	// flag.StringVar(&cfg.db.dsn, "dsn", "{POSTGRES_URI}", "Postgres connection string")
	// flag.StringVar(&cfg.jwt.secret, "jwt-secret", "{JWT_SECRET}", "Secret")
	// flag.Parse()

	cfg := config{
		port: 4000,
		env:  "development",
		db: DB{
			dsn: os.Getenv("POSTGRES_URI"),
		},
		jwt: JWT{
			secret: os.Getenv("JWT_SECRET"),
		},
	}

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModel(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Server running on port", cfg.port)

	logger.Fatalln(srv.ListenAndServe())
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
