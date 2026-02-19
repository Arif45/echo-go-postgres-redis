package routes

import (
	addressRepo "fin-auth/address/repo"
	addressService "fin-auth/address/service"
	authMiddleware "fin-auth/auth/middleware"
	authRes "fin-auth/auth/repo"
	authRest "fin-auth/auth/rest"
	authService "fin-auth/auth/service"
	"fin-auth/cache"
	"fin-auth/config"
	customerRepo "fin-auth/customer/repo"
	customerRest "fin-auth/customer/rest"
	customerService "fin-auth/customer/service"
	personRepo "fin-auth/person/repo"
	personService "fin-auth/person/service"
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
	or := authRes.NewAuthRespository(db, redisCache)
	authSvc := authService.NewAuthService(or, redisCache)

	authRest.SetupAuthRoutes(api, authSvc)

	protected := api.Group("")
	protected.Use(authMiddleware.AuthMiddleware(or, redisCache))
	protected.Use(authMiddleware.RateLimitMiddleware(redisCache, "api"))
	authRest.SetupProtectedRoutes(protected, authSvc)

	cr := customerRepo.NewCustomerRepository(db, redisCache)
	pr := personRepo.NewPersonRepository(db, redisCache)
	ps := personService.NewPersonService(pr, redisCache)
	ar := addressRepo.NewAddressRepository(db, redisCache)
	as := addressService.NewAddressService(ar, redisCache)
	customerSvc := customerService.NewCustomerService(cr, ps, as, redisCache)
	customerRest.SetupCustomerRoutes(protected, customerSvc)
	SetupHealthRoutes(e, db)
}
