package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Vractos/dolly/adapter/api/handler"
	mdw "github.com/Vractos/dolly/adapter/api/middleware"
	"github.com/Vractos/dolly/adapter/cache"
	"github.com/Vractos/dolly/adapter/mercadolivre"
	"github.com/Vractos/dolly/adapter/queue"
	"github.com/Vractos/dolly/adapter/repository"
	"github.com/Vractos/dolly/usecases/announcement"
	"github.com/Vractos/dolly/usecases/order"
	"github.com/Vractos/dolly/usecases/store"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	validate := validator.New()

	// AWS SDK
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Panic("Failed to load config: " + err.Error())
	}

	/// SQS
	client := sqs.NewFromConfig(cfg)

	// PostgreSQL
	dataSourceName := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB_NAME"))
	dbpool, err := pgxpool.New(context.Background(), dataSourceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":6379",
		Password: "",
		DB:       0,
	})
	pong, err := rdb.Ping(rdb.Context()).Result()
	log.Println(pong, err)

	// Order Queue
	orderChan := make(chan []order.OrderMessage)
	orderQueue := queue.NewOrderQueue(client, os.Getenv("ORDER_QUEUE_URL"))

	// Mercado Livre
	mercadoLivre := mercadolivre.NewMercadoLivre(os.Getenv("MELI_APP_ID"), os.Getenv("MELI_SECRET_KEY"), os.Getenv("MELI_REDIRECT_URL"), os.Getenv("MELI_ENDPOINT"), validate)
	// Repositories
	storeRepo := repository.NewStorePostgreSQL(dbpool)
	orderRepo := repository.NewOrderPostgreSQL(dbpool)
	// Caches
	orderCache := cache.NewOrderRedis(rdb)
	// Services
	storeService := store.NewStoreService(storeRepo, mercadoLivre)
	announceService := announcement.NewAnnouncementService(mercadoLivre)
	orderService := order.NewOrderService(
		orderQueue,
		mercadoLivre,
		storeService,
		announceService,
		orderRepo,
		orderCache,
	)

	// Pull messages from queue
	go func() {
		ticker := time.Tick(time.Minute)
		for range ticker {
			log.Println("Pulling...")
			orderChan <- orderQueue.ConsumeOrderNotification()
		}
	}()

	go func() {
		for msgs := range orderChan {
			for _, msg := range msgs {
				orderService.ProcessOrder(msg)
			}
		}
	}()

	// TODO: Make our own router from scratch, based in Radix Tree
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Public Routes
	r.Group(func(r chi.Router) {
		// "/store"
		handler.MakeStoreHandlers(r, storeService)
		handler.MakeOrderHandlers(r, orderService)
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
