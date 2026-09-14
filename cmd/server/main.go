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

	"CineStream/internal/config"
	"CineStream/internal/database"
)

func main() {
	cfg := config.LoadConfig()

	// Database
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

	// Echo
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/health", func(c echo.Context) error {
		if err := db.PingContext(c.Request().Context()); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
				"status":  "UNHEALTHY",
				"message": "Database ping failed",
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "HEALTHY",
			"timestamp": time.Now().Format(time.RFC3339),
			"service":   "CineStream API",
		})
	})

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
