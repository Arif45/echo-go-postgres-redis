package repo

import (
	"context"
	"fin-auth/cache"
	"fin-auth/domain"
	"fin-auth/models"

	"gorm.io/gorm"
)

type Customer struct {
	db    *gorm.DB
	cache *cache.RedisCache
}

func NewCustomerRepository(db *gorm.DB, cache *cache.RedisCache) *Customer {
	return &Customer{
		db:    db,
		cache: cache,
	}
}

func (o *Customer) GetDB(tx ...*gorm.DB) *gorm.DB {
	db := o.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db
}

func (o *Customer) Create(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	if err := o.db.Create(customer).Error; err != nil {
		return nil, err
	}
	return customer, nil
}

func (o *Customer) CreateIndividualCustomer(ctx context.Context, customer *models.Customer, person *models.Person, address *models.Address, req interface{}) (*models.Customer, *models.Person, *models.Address, error) {
	err := o.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(customer).Error; err != nil {
			return err
		}
		person.CustomerID = &customer.ID
		address.CustomerID = &customer.ID

		if err := tx.Create(person).Error; err != nil {
			return err
		}

		if err := tx.Create(address).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, nil, err
	}

	if o.cache != nil && customer.ClientID != nil {
		clientId := *customer.ClientID
		cached := &domain.CachedCustomer{
			Customer: customer,
			Person:   person,
			Address:  address,
		}
		o.cache.CacheCustomer(ctx, clientId, customer.ID, cached)
		o.cache.InvalidateCustomerList(ctx, clientId)
	}

	return customer, person, address, nil
}

func (o *Customer) ListByClientID(ctx context.Context, clientId string) ([]*domain.CachedCustomer, error) {
	var customers []models.Customer
	if err := o.db.Where("client_id = ?", clientId).Find(&customers).Error; err != nil {
		return nil, err
	}

	result := make([]*domain.CachedCustomer, 0, len(customers))
	for i := range customers {
		cached := &domain.CachedCustomer{Customer: &customers[i]}

		var person models.Person
		if err := o.db.Where("customer_id = ?", customers[i].ID).First(&person).Error; err == nil {
			cached.Person = &person
		}

		var addr models.Address
		if err := o.db.Where("customer_id = ?", customers[i].ID).First(&addr).Error; err == nil {
			cached.Address = &addr
		}

		result = append(result, cached)
	}

	return result, nil
}
