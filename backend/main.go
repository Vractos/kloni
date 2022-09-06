package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Vractos/dolly/backend/infrastructure/api/handler"
	"github.com/Vractos/dolly/backend/usecases/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	conn, err := pgx.Connect(context.Background())

	storeService := store.NewStoreService()

	// TODO: Make our own router from scratch, based in Radix Tree
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	handler.MakeStoreHandlers(r, storeService)

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Listing on 8080...")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Panic(err.Error())
	}
}
