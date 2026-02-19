package domain

import (
	"context"
	"fin-auth/models"
)

type AddressService interface {
	Create(ctx context.Context, address *models.Address) (*models.Address, error)
}

type AddressRepository interface {
	Create(ctx context.Context, address *models.Address) (*models.Address, error)
}
