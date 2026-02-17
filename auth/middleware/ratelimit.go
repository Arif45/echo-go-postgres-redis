package middleware

import (
	"fin-auth/cache"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func RateLimitMiddleware(redisCache *cache.RedisCache, limitType string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if redisCache == nil {
				return next(c)
			}
			ip := c.RealIP()
			key := fmt.Sprintf("ratelimit:%s:%s", limitType, ip)
			var max int
			var ttl time.Duration

			switch limitType {
			case "login":
				max = 5
				ttl = 15 * time.Minute
			case "register":
				max = 3
				ttl = 1 * time.Hour
			case "api":
				max = 100
				ttl = 1 * time.Minute
			default:
				max = 100
				ttl = 1 * time.Minute
			}
			allowed, err := redisCache.CheckRateLimit(c.Request().Context(), key, max, ttl)
			if err != nil {
				return next(c)
			}

			if !allowed {
				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"success": false,
					"message": "Rate limit exceeded. Please try again later.",
					"error":   fmt.Sprintf("Maximum %d requests per %v", max, ttl),
				})
			}

			if err := redisCache.IncrementRateLimit(c.Request().Context(), key, ttl); err != nil {
				// Log error but don't block request
			}

			return next(c)
		}
	}
}
