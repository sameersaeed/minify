package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"minify/internal/config"
	"minify/internal/database"
	"minify/internal/handlers"
	"minify/internal/limiter"
	"minify/internal/metrics"
	"minify/internal/middleware"
	"minify/internal/services"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	const maxBuckets = 100000

	// load + validate configs from .env
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	// database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Prometheus metrics
	metrics.Init()

	// services
	urlService := services.NewURLService(db)
	userService := services.NewUserService(db)
	analyticsService := services.NewAnalyticsService(db)
	limiterService := limiter.NewLimiter(maxBuckets)

	// handlers
	urlHandler := handlers.NewURLHandler(urlService, analyticsService, limiterService)
	userHandler := handlers.NewUserHandler(userService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	router := mux.NewRouter()

	router.Use(middleware.CORS(cfg.FrontendURL))
	router.Use(middleware.Logging)
	router.Use(middleware.Metrics)

	setupRoutes(router, urlHandler, userHandler, analyticsHandler)
	router.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// start server in background
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// try to shutdown gracefully on termination (SIGINT/SIGTERM signal)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

// setupRoutes connects handlers to their endpoints
func setupRoutes(router *mux.Router, urlHandler *handlers.URLHandler, userHandler *handlers.UserHandler, analyticsHandler *handlers.AnalyticsHandler) {
	api := router.PathPrefix("/api/v1").Subrouter()

	api.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// URL shortening
	api.HandleFunc("/minify", urlHandler.MinifyURL).Methods("POST")
	api.HandleFunc("/urls", urlHandler.GetUserURLs).Methods("GET")

	// user
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users/login", userHandler.LoginUser).Methods("POST")

	// analytics
	api.HandleFunc("/analytics/overview", analyticsHandler.GetOverview).Methods("GET")
	api.HandleFunc("/analytics/popular", analyticsHandler.GetPopularURLs).Methods("GET")
	api.HandleFunc("/analytics/timeframe/{period}", analyticsHandler.GetTimeframeStats).Methods("GET")

	// redirect + healthcheck
	router.HandleFunc("/{shortCode}", urlHandler.RedirectURL).Methods("GET")
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
}
