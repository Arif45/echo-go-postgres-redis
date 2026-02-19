package service

import (
	"context"
	"fin-auth/cache"
	"fin-auth/domain"
	"fin-auth/models"
)

type Address struct {
	AddressRepository domain.AddressRepository
	Cache             *cache.RedisCache
}

func NewAddressService(repo domain.AddressRepository, cache *cache.RedisCache) *Address {
	return &Address{
		AddressRepository: repo,
		Cache:             cache,
	}
}

func (s *Address) Create(ctx context.Context, address *models.Address) (*models.Address, error) {
	createdAddress, err := s.AddressRepository.Create(ctx, address)
	if err != nil {
		return nil, err
	}
	return createdAddress, nil
}
