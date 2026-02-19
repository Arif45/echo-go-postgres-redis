package repo

import (
	"context"
	"fin-auth/cache"
	"fin-auth/models"

	"gorm.io/gorm"
)

type Address struct {
	db    *gorm.DB
	cache *cache.RedisCache
}

func NewAddressRepository(db *gorm.DB, cache *cache.RedisCache) *Address {
	return &Address{
		db:    db,
		cache: cache,
	}
}

func (o *Address) GetDB(tx ...*gorm.DB) *gorm.DB {
	db := o.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db
}

func (o *Address) Create(ctx context.Context, address *models.Address) (*models.Address, error) {
	if err := o.db.Create(address).Error; err != nil {
		return nil, err
	}
	return address, nil
}
