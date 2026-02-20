package worker

import (
	"fin-auth/models"
	"log"
	"time"

	"gorm.io/gorm"
)

type TokenCleanupWorker struct {
	db *gorm.DB
}

func NewTokenCleanupWorker(db *gorm.DB) *TokenCleanupWorker {
	return &TokenCleanupWorker{
		db: db,
	}
}

func (t *TokenCleanupWorker) Run() error {
	now := time.Now().UTC()
	accessExpiredResult := t.db.Where("expired_at < ?", now).Delete(&models.AccessToken{})

	if accessExpiredResult.Error != nil {
		return accessExpiredResult.Error
	}

	refreshExpiredResult := t.db.Where("expired_at < ?", now).Delete(&models.RefreshToken{})

	if refreshExpiredResult.Error != nil {
		return refreshExpiredResult.Error
	}

	log.Printf("[token-cleanuo] deleted %d access token, %d refresh token", accessExpiredResult.RowsAffected, refreshExpiredResult.RowsAffected)
	return nil
}
