package main

import (
	"hook-api/handlers"
	"hook-api/services"
	"log"
	"net/http"
)

func main() {

	// Connect to DynamoDB
	services.ConnectDynamoDB()

	// Define routes and handlers for sending SMS messages
	http.HandleFunc("/add", handlers.AddJob)
	http.HandleFunc("/test", handlers.Test)

	// Start the HTTP server for the publisher microservice
	log.Println("Server listening on 8080")
	http.ListenAndServe(":8080", nil)
}
