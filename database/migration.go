package database

import (
	"fin-auth/models"
	"log"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database migration...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Secret{},
		&models.AccessToken{},
		&models.RefreshToken{},
		&models.Customer{},
		&models.Person{},
		&models.Address{},
	)

	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}

func GetAllModels() []interface{} {
	return []interface{}{
		&models.User{},
		&models.Secret{},
		&models.AccessToken{},
		&models.RefreshToken{},
		&models.Customer{},
		&models.Person{},
		&models.Address{},
	}
}
