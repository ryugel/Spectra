package main

import (
	"log"
	"net/http"
	"spectra/internal/router"
)

func main() {

	mux := router.NewRouter()

	serverAddr := ":8080"
	log.Printf("Server started at %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
