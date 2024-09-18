package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tablegpt_meter/schemas"
)

// RedisStore implements TokenStore for Redis.
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore creates a new RedisStore.
func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

// GetUserToken retrieves the token information for a user from Redis.
func (s *RedisStore) GetUserToken(ctx context.Context, userID, domain string) (*schemas.UserToken, error) {
	totalTokensKey := fmt.Sprintf("total_tokens:%s:%s", userID, domain)
	usedTokensKey := fmt.Sprintf("used_tokens:%s:%s", userID, domain)

	totalTokensStr, err := s.client.Get(ctx, totalTokensKey).Result()
	if err == redis.Nil {
		return &schemas.UserToken{UserID: userID, TotalTokens: 0, UsedTokens: 0}, nil
	} else if err != nil {
		return nil, err
	}

	totalTokens := 0
	_, err = fmt.Sscanf(totalTokensStr, "%d", &totalTokens)
	if err != nil {
		return nil, err
	}

	usedTokensStr, err := s.client.LRange(ctx, usedTokensKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var usedTokens int32
	for _, record := range usedTokensStr {
		var entry map[string]interface{}
		err = json.Unmarshal([]byte(record), &entry)
		if err != nil {
			return nil, err
		}
		tokens, ok := entry["tokens"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid token record format")
		}
		usedTokens += int32(tokens)
	}

	return &schemas.UserToken{
		UserID:      userID,
		TotalTokens: int32(totalTokens),
		UsedTokens:  usedTokens,
	}, nil
}

// SaveUserToken updates the total token count in Redis.
func (s *RedisStore) SaveUserToken(ctx context.Context, userID, domain string, addedTokens int32) error {
	totalTokensKey := fmt.Sprintf("total_tokens:%s:%s", userID, domain)
	_, err := s.client.IncrBy(ctx, totalTokensKey, int64(addedTokens)).Result()
	return err
}

// SaveTokenUsageRecord records a new token usage entry in Redis.
func (s *RedisStore) SaveTokenUsageRecord(ctx context.Context, userID, domain string, tokensUsed int32) error {
	usedTokensKey := fmt.Sprintf("used_tokens:%s:%s", userID, domain)
	record := map[string]interface{}{
		"tokens":  tokensUsed,
		"used_at": time.Now().UTC().Format(time.RFC3339),
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return s.client.RPush(ctx, usedTokensKey, recordJSON).Err()
}
