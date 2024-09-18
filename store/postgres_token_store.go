package store

import (
	"context"
	"time"

	"github.com/tablegpt_meter/models"
	"github.com/tablegpt_meter/schemas"
	"gorm.io/gorm"
)

// PostgresStore implements TokenStore for PostgreSQL using GORM.
type PostgresStore struct {
	db *gorm.DB
}

// NewPostgresStore creates a new PostgresStore.
func NewPostgresStore(db *gorm.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// GetUserToken retrieves the token information for a user from PostgreSQL.
func (s *PostgresStore) GetUserToken(ctx context.Context, userID, domain string) (*schemas.UserToken, error) {
	var totalToken models.UserTokens
	result := s.db.WithContext(ctx).Where("user_id = ? and domain = ?", userID, domain).First(&totalToken)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &schemas.UserToken{UserID: userID, TotalTokens: 0, UsedTokens: 0}, nil
		}
		return nil, result.Error
	}

	var usedTokens int32
	result = s.db.
		Raw("SELECT COALESCE(SUM(tokens), 0) FROM used_tokens WHERE user_id = ? and domain = ?", userID, domain).
		Scan(&usedTokens)
	if result.Error != nil {
		return nil, result.Error
	}

	return &schemas.UserToken{
		UserID:      totalToken.UserID,
		TotalTokens: totalToken.TotalTokens,
		UsedTokens:  usedTokens,
	}, nil
}

// SaveUserToken updates the total token count in PostgreSQL.
func (s *PostgresStore) SaveUserToken(ctx context.Context, userID, domain string, addedTokens int32) error {
	// Define the UserTokens model
	userTokens := models.UserTokens{
		UserID: userID,
		Domain: domain,
	}

	// Check if the record exists
	err := s.db.WithContext(ctx).Where("user_id = ? AND domain = ?", userID, domain).
		First(&userTokens).Error

	if err == gorm.ErrRecordNotFound {
		// Record not found, create a new one with the initial token value
		userTokens.TotalTokens = addedTokens
		return s.db.WithContext(ctx).Create(&userTokens).Error
	} else if err != nil {
		// Return any other errors that occurred during the find operation
		return err
	}

	// Record exists, update the total_tokens field
	return s.db.WithContext(ctx).Model(&userTokens).
		Update("total_tokens", gorm.Expr("total_tokens + ?", addedTokens)).Error
}

// SaveTokenUsageRecord records a new token usage entry.
func (s *PostgresStore) SaveTokenUsageRecord(ctx context.Context, userID, domain string, tokensUsed int32) error {
	usedToken := models.UsedToken{
		UserID: userID,
		Domain: domain,
		Tokens: tokensUsed,
		UsedAt: time.Now(),
	}
	return s.db.WithContext(ctx).Create(&usedToken).Error
}
