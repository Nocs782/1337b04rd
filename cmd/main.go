package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"1337b04rd/internal/adapter/postgres"
	"1337b04rd/internal/handler"
	"1337b04rd/internal/s3"
)

func main() {

	time.Sleep(10 * time.Second)

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

	s3Endpoint := os.Getenv("S3_ENDPOINT")
	s3Bucket := os.Getenv("S3_BUCKET")

	if s3Bucket == "" {
		log.Fatal("S3_BUCKET is not set.")
	}

	minioStorage, err := s3.NewMinioStorage(s3Endpoint, s3Bucket) // FIXED: catch err here!
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	fmt.Println("S3 client initialized with bucket:", s3Bucket)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, db, minioStorage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server is running on port %s\n", port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
