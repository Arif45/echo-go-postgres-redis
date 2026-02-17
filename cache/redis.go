package cache

import (
	"context"
	"encoding/json"
	"fin-auth/domain"
	"fin-auth/models"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (r *RedisCache) CacheAccessToken(ctx context.Context, token string, data *models.AccessToken) error {
	key := fmt.Sprintf("token:access:%s", token)

	tokenData := map[string]interface{}{
		"client_id":  data.ClientId,
		"expired_at": data.ExpiredAt.Unix(),
		"created_at": data.CreatedAt.Unix(),
	}

	ttl := time.Until(data.ExpiredAt)
	if ttl <= 0 {
		return nil
	}

	return r.client.HSet(ctx, key, tokenData).Err()
}

func (r *RedisCache) CacheRefreshToken(ctx context.Context, token string, data *models.RefreshToken) error {
	key := fmt.Sprintf("token:refresh:%s", token)

	tokenData := map[string]interface{}{
		"client_id":       data.ClientId,
		"access_token_id": data.AccessTokenId,
		"expired_at":      data.ExpiredAt.Unix(),
		"created_at":      data.CreatedAt.Unix(),
	}

	ttl := time.Until(data.ExpiredAt)
	if ttl <= 0 {
		return nil
	}

	return r.client.HSet(ctx, key, tokenData).Err()
}

func (r *RedisCache) GetAccessToken(ctx context.Context, token string) (*models.AccessToken, error) {
	key := fmt.Sprintf("token:access:%s", token)

	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("token not found in cache")
	}

	expiredAt, _ := strconv.ParseInt(result["expired_at"], 10, 64)
	createdAt, _ := strconv.ParseInt(result["created_at"], 10, 64)

	accessToken := &models.AccessToken{
		ClientId:  result["client_id"],
		Token:     token,
		ExpiredAt: time.Unix(expiredAt, 0),
	}
	accessToken.CreatedAt = time.Unix(createdAt, 0)

	if accessToken.ExpiredAt.Before(time.Now().UTC()) {
		r.client.Del(ctx, key)
		return nil, fmt.Errorf("token expired")
	}

	return accessToken, nil
}

func (r *RedisCache) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	key := fmt.Sprintf("token:refresh:%s", token)

	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("token not found in cache")
	}

	expiredAt, _ := strconv.ParseInt(result["expired_at"], 10, 64)
	createdAt, _ := strconv.ParseInt(result["created_at"], 10, 64)
	accessTokenId, _ := strconv.ParseUint(result["access_token_id"], 10, 32)

	refreshToken := &models.RefreshToken{
		ClientId:      result["client_id"],
		Token:         token,
		AccessTokenId: uint(accessTokenId),
		ExpiredAt:     time.Unix(expiredAt, 0),
	}
	refreshToken.CreatedAt = time.Unix(createdAt, 0)

	if refreshToken.ExpiredAt.Before(time.Now().UTC()) {
		r.client.Del(ctx, key)
		return nil, fmt.Errorf("token expired")
	}

	return refreshToken, nil
}

// Blacklist Operations

func (r *RedisCache) BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	key := fmt.Sprintf("blacklist:token:%s", token)
	return r.client.Set(ctx, key, "1", ttl).Err()
}

func (r *RedisCache) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:token:%s", token)
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *RedisCache) DeleteTokenFromCache(ctx context.Context, token string) error {
	accessKey := fmt.Sprintf("token:access:%s", token)
	return r.client.Del(ctx, accessKey).Err()
}

// Session Operations

func (r *RedisCache) CreateSession(ctx context.Context, clientId, token string, data *models.SessionData) error {
	key := fmt.Sprintf("session:%s:%s", clientId, token)

	sessionJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, sessionJSON, 24*time.Hour).Err()
}

func (r *RedisCache) GetSession(ctx context.Context, clientId, token string) (*models.SessionData, error) {
	key := fmt.Sprintf("session:%s:%s", clientId, token)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var session models.SessionData
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *RedisCache) UpdateLastActivity(ctx context.Context, clientId, token string) error {
	key := fmt.Sprintf("session:%s:%s", clientId, token)

	session, err := r.GetSession(ctx, clientId, token)
	if err != nil {
		return nil
	}
	session.LastActivity = time.Now().UTC()

	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// Reset TTL
	return r.client.Set(ctx, key, sessionJSON, 24*time.Hour).Err()
}

func (r *RedisCache) ListSessions(ctx context.Context, clientId string) ([]*models.SessionData, error) {
	pattern := fmt.Sprintf("session:%s:*", clientId)

	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	sessions := make([]*models.SessionData, 0, len(keys))
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue // Skip failed reads
		}

		var session models.SessionData
		if err := json.Unmarshal([]byte(data), &session); err != nil {
			continue
		}

		sessions = append(sessions, &session)
	}

	return sessions, nil
}

func (r *RedisCache) DeleteSession(ctx context.Context, clientId, token string) error {
	key := fmt.Sprintf("session:%s:%s", clientId, token)
	return r.client.Del(ctx, key).Err()
}

// Rate Limiting Operations

func (r *RedisCache) CheckRateLimit(ctx context.Context, key string, max int, ttl time.Duration) (bool, error) {
	count, err := r.client.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return true, err // Allow on error
	}

	return count < max, nil
}

func (r *RedisCache) IncrementRateLimit(ctx context.Context, key string, ttl time.Duration) error {
	pipe := r.client.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, ttl)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *RedisCache) GetRateLimitCount(ctx context.Context, key string) (int, error) {
	count, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return count, err
}

func (r *RedisCache) CacheClient(ctx context.Context, clientId string, data *domain.ClientWithSecrets) error {
	key := fmt.Sprintf("client:%s", clientId)

	clientData := map[string]interface{}{
		"name":             data.User.Name,
		"email":            data.User.Email,
		"is_active":        data.User.IsActive,
		"secret":           data.Secret.Secret,
		"secondary_secret": data.Secret.SecondarySecret,
	}

	err := r.client.HSet(ctx, key, clientData).Err()
	if err != nil {
		return err
	}

	return r.client.Expire(ctx, key, 5*time.Minute).Err()
}

func (r *RedisCache) GetCachedClient(ctx context.Context, clientId string) (*domain.ClientWithSecrets, error) {
	key := fmt.Sprintf("client:%s", clientId)

	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("client not found in cache")
	}

	isActive, _ := strconv.ParseBool(result["is_active"])

	user := &models.User{
		ClientId: clientId,
		Name:     result["name"],
		Email:    result["email"],
		IsActive: isActive,
	}

	secret := &models.Secret{
		ClientId:        clientId,
		Secret:          result["secret"],
		SecondarySecret: result["secondary_secret"],
	}

	return &domain.ClientWithSecrets{
		User:   user,
		Secret: secret,
	}, nil
}

func (r *RedisCache) InvalidateClient(ctx context.Context, clientId string) error {
	key := fmt.Sprintf("client:%s", clientId)
	return r.client.Del(ctx, key).Err()
}

// Health Check

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
