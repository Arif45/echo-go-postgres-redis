package service

import (
	"context"
	"errors"
	"fin-auth/cache"
	"fin-auth/domain"
	"fin-auth/models"
)

type Customer struct {
	CustomerRepository domain.CustomerRepository
	PersonService      domain.PersonService
	AddressService     domain.AddressService
	Cache              *cache.RedisCache
}

func NewCustomerService(repo domain.CustomerRepository, personService domain.PersonService, addressService domain.AddressService, cache *cache.RedisCache) *Customer {
	return &Customer{
		CustomerRepository: repo,
		PersonService:      personService,
		AddressService:     addressService,
		Cache:              cache,
	}
}

func (s *Customer) CreateIndividualCustomer(ctx context.Context, customer *models.Customer, person *models.Person, address *models.Address, req interface{}) (*models.Customer, *models.Person, *models.Address, error) {
	if customer == nil {
		return nil, nil, nil, errors.New("Invalid customer data")
	}
	if person == nil {
		return nil, nil, nil, errors.New("Invalid person data")
	}
	if address == nil {
		return nil, nil, nil, errors.New("Invalid address data")
	}
	createdCustomer, createdPerson, createdAddress, err := s.CustomerRepository.CreateIndividualCustomer(ctx, customer, person, address, req)
	if err != nil {
		return nil, nil, nil, err
	}
	return createdCustomer, createdPerson, createdAddress, nil
}

func (s *Customer) GetCustomersByClientID(ctx context.Context, clientId string) ([]*domain.CachedCustomer, error) {
	if s.Cache != nil {
		if customers, err := s.Cache.GetCachedCustomerList(ctx, clientId); err == nil {
			return customers, nil
		}
	}

	customers, err := s.CustomerRepository.ListByClientID(ctx, clientId)
	if err != nil {
		return nil, err
	}

	if s.Cache != nil {
		s.Cache.CacheCustomerList(ctx, clientId, customers)
	}

	return customers, nil
}
