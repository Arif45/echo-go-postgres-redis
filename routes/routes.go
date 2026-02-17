package routes

import (
	authMiddleware "fin-auth/auth/middleware"
	authRes "fin-auth/auth/repo"
	authRest "fin-auth/auth/rest"
	authService "fin-auth/auth/service"
	"fin-auth/cache"
	"fin-auth/config"
	"log"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	db, err := config.InitGormDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	redisClient, err := config.InitRedis()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v (continuing without cache)", err)
		redisClient = nil
	}

	var redisCache *cache.RedisCache
	if redisClient != nil {
		redisCache = cache.NewRedisCache(redisClient)
	}

	api := e.Group("/api/v1")
	or := authRes.New(db, redisCache)
	authSvc := authService.New(or, redisCache)

	authRest.SetupAuthRoutes(api, authSvc)
	protected := api.Group("")
	protected.Use(authMiddleware.AuthMiddleware(or, redisCache))
	protected.Use(authMiddleware.RateLimitMiddleware(redisCache, "api"))
	authRest.SetupProtectedRoutes(protected, authSvc)

	SetupHealthRoutes(e, db)
}
