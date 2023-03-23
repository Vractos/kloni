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
	"github.com/Vractos/dolly/pkg/metrics"
	"github.com/Vractos/dolly/usecases/announcement"
	"github.com/Vractos/dolly/usecases/order"
	"github.com/Vractos/dolly/usecases/store"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func startEnv() {
	if env := os.Getenv("APP_ENV"); env == "" || env == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}

func main() {
	// ENV
	startEnv()

	// Log
	logger := metrics.NewLogger("info")
	defer logger.Sync()
	// Tracer

	// Validator package
	validate := validator.New()

	// AWS SDK
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		logger.Panic("Failed to load config: "+err.Error(), err)
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
	logger.Warn(pong,
		zap.Error(err),
	)

	// Order Queue
	orderChan := make(chan []order.OrderMessage)
	orderQueue := queue.NewOrderQueue(client, os.Getenv("ORDER_QUEUE_URL"), *logger)

	// Mercado Livre
	mercadoLivre := mercadolivre.NewMercadoLivre(os.Getenv("MELI_APP_ID"), os.Getenv("MELI_SECRET_KEY"), os.Getenv("MELI_REDIRECT_URL"), os.Getenv("MELI_ENDPOINT"), validate, *logger)
	// Repositories
	storeRepo := repository.NewStorePostgreSQL(dbpool, *logger)
	orderRepo := repository.NewOrderPostgreSQL(dbpool, *logger)
	// Caches
	orderCache := cache.NewOrderRedis(rdb)
	// Services
	storeService := store.NewStoreService(storeRepo, mercadoLivre, *logger)
	announceService := announcement.NewAnnouncementService(mercadoLivre, storeService, *logger)
	orderService := order.NewOrderService(
		orderQueue,
		mercadoLivre,
		storeService,
		announceService,
		orderRepo,
		orderCache,
		*logger,
	)

	// Pull messages from queue
	go func() {
		ticker := time.Tick(time.Minute)
		for range ticker {
			orderChan <- orderQueue.ConsumeOrderNotification()
		}
	}()

	// Process messages
	go func() {
		for msgs := range orderChan {
			for _, msg := range msgs {
				orderService.ProcessOrder(msg)
			}
		}
	}()

	// Router
	// TODO Make our own router from scratch, based in Radix Tree
	r := chi.NewRouter()
	r.Use(mdw.NewStructuredLogger(logger))
	// r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Public Routes
	r.Group(func(r chi.Router) {
		// "/store"
		handler.MakeStoreHandlers(r, storeService, *logger)
		// "/order"
		handler.MakeOrderHandlers(r, orderService, *logger)
	})

	// Private Routes
	r.Group(func(r chi.Router) {
		r.Use(mdw.EnsureValidToken(*logger))
		r.Use(mdw.AddStoreIDToCtx)

		handler.MakeAnnouncementHandlers(r, announceService, storeService, *logger)
	})

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger.Info("Listing on 80")
	err = http.ListenAndServe(":80", r)
	if err != nil {
		logger.Panic(err.Error(), err)
	}
}
