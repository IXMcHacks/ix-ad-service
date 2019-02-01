package main

import (
	"log"
	"net/http"
	"strconv"

	"ixmchacks/ix-ad-service/handlers"
)

const (
	port = 8080
)

func main() {

	// Use Log package to make print statements to standard output
	log.Printf("Starting up ad-server listening on port: %v", port)

	// 1. Register a HTTP Request Handler
	http.HandleFunc("/ixrtb", handlers.RunAuction)

	// 2. Call Listen and Server and provide the ipaddress and port you want
	// to reach the server on. Adding just the port automatically sets the server
	// listening on localhost, or 127.0.0.1
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)

	if err != nil {
		log.Printf("Error starting up server:%", err)
	}
}
