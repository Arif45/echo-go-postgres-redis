package repo

import (
	"context"
	"fin-auth/cache"
	"fin-auth/models"

	"gorm.io/gorm"
)

type Person struct {
	db    *gorm.DB
	cache *cache.RedisCache
}

func NewPersonRepository(db *gorm.DB, cache *cache.RedisCache) *Person {
	return &Person{
		db:    db,
		cache: cache,
	}
}

func (o *Person) GetDB(tx ...*gorm.DB) *gorm.DB {
	db := o.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db
}

func (o *Person) Create(ctx context.Context, person *models.Person) (*models.Person, error) {
	if err := o.db.Create(person).Error; err != nil {
		return nil, err
	}
	return person, nil
}
