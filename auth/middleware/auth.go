package middleware

import (
	"fin-auth/cache"
	"fin-auth/domain"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(repo domain.AuthRepository, redisCache *cache.RedisCache) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization format")
			}

			token := parts[1]
			if redisCache != nil {
				blacklisted, err := redisCache.IsTokenBlacklisted(c.Request().Context(), token)
				if err == nil && blacklisted {
					return echo.NewHTTPError(http.StatusUnauthorized, "User has been logged out")
				}
			}
			accessToken, err := repo.FindValidAccessToken(c.Request().Context(), token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
			}
			if redisCache != nil {
				redisCache.UpdateLastActivity(c.Request().Context(), accessToken.ClientId, token) // Ignore errors
			}
			c.Set("client_id", accessToken.ClientId)
			c.Set("token", token)

			return next(c)
		}
	}
}
