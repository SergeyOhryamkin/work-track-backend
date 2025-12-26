package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/sergey/work-track-backend/internal/config"
	"github.com/sergey/work-track-backend/internal/database"
	"github.com/sergey/work-track-backend/internal/handler"
	"github.com/sergey/work-track-backend/internal/middleware"
	"github.com/sergey/work-track-backend/internal/repository"
	"github.com/sergey/work-track-backend/internal/service"
)

func main() {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.NewSQLiteDB(cfg.Database.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	log.Println("Successfully connected to database")

	// Run migrations
	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	trackItemRepo := repository.NewTrackItemRepository(db)
	sessionRepo := repository.NewUserSessionRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, sessionRepo, cfg.JWT.Secret)
	trackItemService := service.NewTrackItemService(trackItemRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	trackItemHandler := handler.NewTrackItemHandler(trackItemService)

	// Setup router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			// Public routes
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)

			// Protected routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
				r.Post("/logout", authHandler.Logout)
			})
		})

		// Track item routes (protected)
		r.Route("/track-items", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
			r.Get("/", trackItemHandler.ListTrackItems)
			r.Post("/", trackItemHandler.CreateTrackItem)
			r.Get("/{id}", trackItemHandler.GetTrackItem)
			r.Put("/{id}", trackItemHandler.UpdateTrackItem)
			r.Delete("/{id}", trackItemHandler.DeleteTrackItem)
		})
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s (environment: %s)", cfg.Server.Port, cfg.Server.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server is shutting down...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
