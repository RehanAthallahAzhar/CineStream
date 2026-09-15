package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"CineStream/internal/auth"
	"CineStream/internal/config"
	"CineStream/internal/database"
	"CineStream/internal/response"
	"CineStream/internal/showtime"
)

func main() {
	cfg := config.LoadConfig()

	// db
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed gracefully")
		}
	}()

	// module
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg)
	authHandler := auth.NewHandler(authService)

	showtimeRepo := showtime.NewRepository(db)
	showtimeService := showtime.NewService(showtimeRepo)
	showtimeHandler := showtime.NewHandler(showtimeService)

	// echo
	e := echo.New()
	e.HideBanner = true

	// custom error handler
	e.HTTPErrorHandler = response.CustomHTTPErrorHandler

	// middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// health check
	e.GET("/health", func(c echo.Context) error {
		if err := db.PingContext(c.Request().Context()); err != nil {
			return response.Error(c, http.StatusServiceUnavailable, "Database ping failed")
		}
		return response.Success(c, http.StatusOK, "Service is healthy", map[string]interface{}{
			"status":    "HEALTHY",
			"timestamp": time.Now().Format(time.RFC3339),
			"service":   "CineStream API",
		})
	})

	// routes
	auth.RegisterRoutes(e, authHandler)
	showtime.RegisterRoutes(e, showtimeHandler, cfg)

	// server start
	serverAddress := fmt.Sprintf(":%s", cfg.AppPort)
	go func() {
		log.Printf("Starting CineStream API server on port %s (%s mode)", cfg.AppPort, cfg.AppEnv)
		if err := e.Start(serverAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Shutting down the server unexpectedly: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received. Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting complete")
}
