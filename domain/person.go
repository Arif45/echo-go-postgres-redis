package domain

import (
	"context"
	"fin-auth/models"
)

type PersonService interface {
	Create(ctx context.Context, person *models.Person) (*models.Person, error)
}

type PersonRepository interface {
	Create(ctx context.Context, person *models.Person) (*models.Person, error)
}
