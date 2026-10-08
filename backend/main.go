package main

import (
	"log"
	"net/http"
	"os"
	"qare/backend/internal/handlers"
	"qare/backend/internal/repositories"
	"qare/backend/internal/routes"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("MYSQL_DSN is required")
	}
	db, err := repositories.Open(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	port := env("PORT", "8080")
	log.Println("API on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, routes.New(handlers.New(db), env("FRONTEND_ORIGIN", "http://localhost:3000"))))
}
