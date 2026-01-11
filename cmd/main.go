package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/suyashmohan/godoit/gen/database"
	"github.com/suyashmohan/godoit/gen/todo/v1/todov1connect"
	"github.com/suyashmohan/godoit/internal/services"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func ConnectionString() string {
	Host := os.Getenv("DB_HOST")
	Port := os.Getenv("DB_PORT")
	User := os.Getenv("DB_USER")
	Password := os.Getenv("DB_PASSWORD")
	DBName := os.Getenv("DB_NAME")
	SSLMode := os.Getenv("DB_SSLMODE")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		Host, Port, User, Password, DBName, SSLMode,
	)
}

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file", err)
	}

	dbConn, err := pgxpool.New(ctx, ConnectionString())
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}
	defer dbConn.Close()

	mux := http.NewServeMux()

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	todoService := services.NewTodoService(database.New(dbConn))
	path, handler := todov1connect.NewTodoServiceHandler(todoService)
	mux.Handle(path, corsMiddleware(handler))

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", h2c.NewHandler(mux, &http2.Server{}))
}
