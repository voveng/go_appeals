package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_appeals/internal/handlers"
	"go_appeals/internal/repository"
	"go_appeals/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())

	dbPath := "./appeals.db"
	repo, err := repository.NewAppealRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	defer func() {
		if err := repo.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
		log.Println("Database connection closed.")
	}()

	service := services.NewAppealService(repo)

	apiHandlers := &handlers.Handlers{
		Service: service,
	}

	api := app.Group("/appeals")
	api.Get("/", apiHandlers.GetStartedAppeals)
	api.Get("/all", apiHandlers.GetAllAppeals)
	api.Get("/by-dates", apiHandlers.GetAppealsByDates)
	api.Post("/cancel-all-in-progress", apiHandlers.CancelAllInProgress)
	api.Get("/:id", apiHandlers.GetAppealByID)
	api.Post("/", apiHandlers.CreateAppeal)
	api.Patch("/:id/start", apiHandlers.StartProcessing)
	api.Patch("/:id/complete", apiHandlers.CompleteAppeal)
	api.Patch("/:id/cancel", apiHandlers.CancelAppeal)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server is starting on port %s...", port)
		if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server gracefully stopped.")
}
