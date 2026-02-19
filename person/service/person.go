package service

import (
	"context"
	"fin-auth/cache"
	"fin-auth/domain"
	"fin-auth/models"
)

type Person struct {
	PersonRepository domain.PersonRepository
	Cache            *cache.RedisCache
}

func NewPersonService(repo domain.PersonRepository, cache *cache.RedisCache) *Person {
	return &Person{
		PersonRepository: repo,
		Cache:            cache,
	}
}

func (s *Person) Create(ctx context.Context, person *models.Person) (*models.Person, error) {
	createdPerson, err := s.PersonRepository.Create(ctx, person)
	if err != nil {
		return nil, err
	}
	return createdPerson, nil
}
