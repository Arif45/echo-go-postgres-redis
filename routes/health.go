package routes

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupHealthRoutes(e *echo.Echo, db *gorm.DB) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to fin-auth API",
			"version": "v1.0.0",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return healthCheck(c, db)
	})
}

func healthCheck(c echo.Context, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed to get database instance: %v", err)
		return c.JSON(http.StatusOK, map[string]string{
			"status":   "healthy",
			"database": "unhealthy",
		})
	}

	dbStatus := "healthy"
	if err := sqlDB.Ping(); err != nil {
		dbStatus = "unhealthy"
		log.Printf("Database health check failed: %v", err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":   "healthy",
		"database": dbStatus,
	})
}
