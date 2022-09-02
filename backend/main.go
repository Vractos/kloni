package main

import (
	"log"
	"net/http"

	"github.com/Vractos/dolly/backend/infrastructure/api/handler"
	"github.com/Vractos/dolly/backend/usecases/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	storeService := store.NewStoreService()

	mux := http.NewServeMux()
	mux.Handle("/store", handler.RegisterStore(storeService))

	log.Println("Listing on 8080...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Panic(err.Error())
	}
}
