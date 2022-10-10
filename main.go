package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Vractos/dolly/adapter/api/handler"
	mdw "github.com/Vractos/dolly/adapter/api/middleware"
	"github.com/Vractos/dolly/adapter/mercadolivre"
	"github.com/Vractos/dolly/adapter/repository"
	"github.com/Vractos/dolly/usecases/store"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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
	// AWS SDK
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Panic("Failed to load config: " + err.Error())
	}
	/// SQS
	sqs := sqs.NewFromConfig(cfg)
	// PostgreSQL
	dataSourceName := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB_NAME"))
	conn, err := pgx.Connect(context.Background(), dataSourceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Mercado Livre
	meliStore := mercadolivre.NewMercadoLivreStore(os.Getenv("MELI_APP_ID"), os.Getenv("MELI_SECRET_KEY"), os.Getenv("MELI_REDIRECT_URL"), os.Getenv("MELI_ENDPOINT"))
	// Repositories
	storeRepo := repository.NewStorePostgreSQL(conn)
	// Services
	storeService := store.NewStoreService(storeRepo, meliStore)

	// TODO: Make our own router from scratch, based in Radix Tree
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Public Routes
	r.Group(func(r chi.Router) {
		// "/store"
		handler.MakeStoreHandlers(r, storeService)
	})

	// Private Routes
	r.Group(func(r chi.Router) {
		r.Use(mdw.EnsureValidToken())
		r.Use(mdw.AddStoreIDToCtx)
	})

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Listing on 8080...")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Panic(err.Error())
	}
}
