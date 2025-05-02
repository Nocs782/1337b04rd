package main

import (
	"fmt"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DATABASE_URL")
	s3Endpoint := os.Getenv("S3_ENDPOINT")

	if port == "" || dbURL == "" || s3Endpoint == "" {
		fmt.Println("One or more environment variables are missing.")
		os.Exit(1)
	}

	fmt.Println("PORT:", port)
	fmt.Println("DATABASE_URL:", dbURL)
	fmt.Println("S3_ENDPOINT:", s3Endpoint)

	fmt.Println("1337b04rd is running!")
}
