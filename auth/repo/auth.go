package repo

import (
	"context"
	"fin-auth/cache"
	"fin-auth/domain"
	"fin-auth/dto"
	"fin-auth/models"
	"time"

	"gorm.io/gorm"
)

type Auth struct {
	db    *gorm.DB
	cache *cache.RedisCache
}

func New(db *gorm.DB, cache *cache.RedisCache) *Auth {
	return &Auth{
		db:    db,
		cache: cache,
	}
}

func (o *Auth) GetDB(tx ...*gorm.DB) *gorm.DB {
	db := o.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db
}

func (auth *Auth) CreateAuthClient(ctx context.Context, user *models.User, secret *models.Secret) (*dto.RegisterClientRes, error) {
	err := auth.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		if err := tx.Create(secret).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	res := &dto.RegisterClientRes{
		ClientId:        user.ClientId,
		Secret:          secret.Secret,
		SecondarySecret: secret.SecondarySecret,
	}

	return res, nil
}

func (auth *Auth) FindClientWithSecrets(ctx context.Context, clientId string) (*domain.ClientWithSecrets, error) {
	if auth.cache != nil {
		cachedClient, err := auth.cache.GetCachedClient(ctx, clientId)
		if err == nil && cachedClient != nil {
			return cachedClient, nil
		}
	}
	var user models.User
	err := auth.db.Where("client_id = ?", clientId).First(&user).Error
	if err != nil {
		return nil, err
	}

	var secret models.Secret
	err = auth.db.Where("client_id = ?", clientId).First(&secret).Error
	if err != nil {
		return nil, err
	}

	result := &domain.ClientWithSecrets{
		User:   &user,
		Secret: &secret,
	}
	if auth.cache != nil {
		auth.cache.CacheClient(ctx, clientId, result)
	}

	return result, nil
}

func (auth *Auth) CreateAccessToken(ctx context.Context, token *models.AccessToken) error {
	if err := auth.db.Create(token).Error; err != nil {
		return err
	}

	if auth.cache != nil {
		auth.cache.CacheAccessToken(ctx, token.Token, token)
	}

	return nil
}

func (auth *Auth) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	if err := auth.db.Create(token).Error; err != nil {
		return err
	}
	if auth.cache != nil {
		auth.cache.CacheRefreshToken(ctx, token.Token, token)
	}

	return nil
}

func (auth *Auth) FindValidRefreshToken(ctx context.Context, tokenStr string) (*models.RefreshToken, error) {
	if auth.cache != nil {
		cachedToken, err := auth.cache.GetRefreshToken(ctx, tokenStr)
		if err == nil && cachedToken != nil {
			return cachedToken, nil
		}
	}
	var token models.RefreshToken
	err := auth.db.Where("token = ? AND expired_at > ?", tokenStr, time.Now().UTC()).First(&token).Error
	if err != nil {
		return nil, err
	}
	if auth.cache != nil {
		auth.cache.CacheRefreshToken(ctx, token.Token, &token)
	}

	return &token, nil
}

func (auth *Auth) FindValidAccessToken(ctx context.Context, tokenStr string) (*models.AccessToken, error) {
	if auth.cache != nil {
		cachedToken, err := auth.cache.GetAccessToken(ctx, tokenStr)
		if err == nil && cachedToken != nil {
			return cachedToken, nil
		}
	}
	var token models.AccessToken
	err := auth.db.Where("token = ? AND expired_at > ?", tokenStr, time.Now().UTC()).First(&token).Error
	if err != nil {
		return nil, err
	}
	if auth.cache != nil {
		auth.cache.CacheAccessToken(ctx, token.Token, &token)
	}

	return &token, nil
}
