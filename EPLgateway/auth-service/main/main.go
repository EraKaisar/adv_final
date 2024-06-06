package main

import (
	data "EPLgateway/auth-service/internal/model"
	"EPLgateway/auth-service/jsonlog"
	"EPLgateway/auth-service/mailer"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	db *sql.DB
)

type config struct {
	port int
	env  string
	db   struct {
		url          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
}
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}
type logger struct {
	out      io.Writer
	minLevel jsonlog.Level
	mu       sync.Mutex
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbURL := os.Getenv("DB_URL")

	db = initDB(dbURL)
	app := &application{
		config: config{
			port: 8080,
			env:  "development",
		}}
	// Applying migrations

	r := app.setupRoutes()
	log.Println("Server started on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func (app *application) setupRoutes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return router
}

func initDB(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	return db
}
