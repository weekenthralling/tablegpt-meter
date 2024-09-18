package store

import (
	"context"

	"github.com/tablegpt_meter/schemas"
)

// TokenStore defines methods for token operations.
type TokenStore interface {
	// GetUserToken retrieves the token information for a user.
	GetUserToken(ctx context.Context, userID, domain string) (*schemas.UserToken, error)
	// SaveUserToken updates the token information for a user.
	SaveUserToken(ctx context.Context, userID, domain string, usedTokens int32) error
	// SaveTokenUsageRecord records a new token usage entry.
	SaveTokenUsageRecord(ctx context.Context, userID, domain string, tokensUsed int32) error
}
