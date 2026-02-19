package domain

import (
	"context"
	"fin-auth/models"
)

type CachedCustomer struct {
	Customer *models.Customer `json:"customer"`
	Person   *models.Person   `json:"person"`
	Address  *models.Address  `json:"address"`
}

type CustomerService interface {
	CreateIndividualCustomer(ctx context.Context, customer *models.Customer, person *models.Person, address *models.Address, req interface{}) (*models.Customer, *models.Person, *models.Address, error)
	GetCustomersByClientID(ctx context.Context, clientId string) ([]*CachedCustomer, error)
}

type CustomerRepository interface {
	Create(ctx context.Context, customer *models.Customer) (*models.Customer, error)
	CreateIndividualCustomer(ctx context.Context, customer *models.Customer, person *models.Person, address *models.Address, req interface{}) (*models.Customer, *models.Person, *models.Address, error)
	ListByClientID(ctx context.Context, clientId string) ([]*CachedCustomer, error)
}
