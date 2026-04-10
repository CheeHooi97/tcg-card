package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"pkm/config"
	"pkm/database"
	"pkm/handler"
	"pkm/repository"
	"pkm/router"
	service "pkm/service"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	// Load config
	config.LoadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Auto-migrate
	err = database.Migrate(db)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repository
	repos := repository.InitializeRepository(db)

	// Initialize services
	services := service.InitializeService(repos)

	// Initialize message handler
	h := handler.NewHandler(services)

	// Setup API routes
	api := router.SetupRoutes(h, db)

	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if !c.Response().Committed {
			c.JSON(http.StatusInternalServerError, map[string]any{
				"error": map[string]any{
					"code":    "INTERNAL_ERROR",
					"message": "Internal error",
					"debug":   err.Error(),
				},
			})
		}
	}

	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		api.ServeHTTP(res, req)
		return
	})
	if err := e.Start(":2001"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
