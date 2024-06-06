package main

import (
	"EPLgateway/auth-service/jsonlog"
	"EPLgateway/comment-service/internal/model"
	"database/sql"
	_ "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
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
	cors struct {
		trustedOrigins []string
	}
	jwt struct {
		secret string
	}
}

type application struct {
	logger *jsonlog.Logger
	config config
	models model.Models
	wg     sync.WaitGroup
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbURL := os.Getenv("DB_URL")

	jwtSecret := os.Getenv("JWT_SECRET")

	db = initDB(dbURL)
	defer db.Close()

	app := &application{
		config: config{
			port: 8081,
			env:  "development",
			jwt:  struct{ secret string }{secret: jwtSecret},
		},
		models: model.NewModels(db),
	}

	r := app.setupRoutes()

	log.Println("Server started on :8081")
	err = http.ListenAndServe(":8081", r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
func (app *application) setupRoutes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/teams/:id/comments", app.requireAuthentication(app.createCommentHandler))
	router.HandlerFunc(http.MethodPost, "/v1/teams/:id/comments", app.requireAuthentication(app.updateCommentHandler))
	router.HandlerFunc(http.MethodPost, "/v1/teams/:id/comments", app.requireAuthentication(app.deleteCommentHandler))
	router.HandlerFunc(http.MethodGet, "/v1/teams/:id/comments", app.listCommentsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/teams/:id/ratings", app.requireAuthentication(app.createRatingHandler))
	router.HandlerFunc(http.MethodPost, "/v1/teams/:id/ratings", app.requireAuthentication(app.updateRatingHandler))
	router.HandlerFunc(http.MethodPost, "/v1/teams/:id/ratings", app.requireAuthentication(app.deleteRatingHandler))
	router.HandlerFunc(http.MethodGet, "/v1/teams/:id/ratings", app.listRatingsHandler)

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
