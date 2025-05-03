package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"1337b04rd/internal/adapter/postgres"
	"1337b04rd/internal/handler"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set.")
	}

	db, err := postgres.Connect(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to the database.")

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, db)

	port := os.Getenv("PORT")

	fmt.Printf("Server is running on port %s\n", port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
