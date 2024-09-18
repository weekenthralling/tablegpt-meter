package server

import (
	"context"
	"log"

	ratelimit "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"
	"github.com/tablegpt_meter/store"
)

type TokenLimitServer struct {
	ratelimit.RateLimitServiceServer
	tokenStore store.TokenStore
}

// ShouldRateLimit implements ratelimitv3.RateLimitServiceServer.
func (t *TokenLimitServer) ShouldRateLimit(ctx context.Context, req *ratelimit.RateLimitRequest) (*ratelimit.RateLimitResponse, error) {
	// Check if the user has exceeded the total number of tokens.
	for _, descriptor := range req.Descriptors {
		userId := descriptor.Entries[0].Value
		userToken, _ := t.tokenStore.GetUserToken(ctx, userId, req.Domain)
		log.Printf("User %s has %d tokens remaining", userId, userToken.TotalTokens-userToken.UsedTokens)
		if userToken.UsedTokens > userToken.TotalTokens {
			return &ratelimit.RateLimitResponse{
				OverallCode: ratelimit.RateLimitResponse_OVER_LIMIT,
			}, nil
		}
	}

	return &ratelimit.RateLimitResponse{
		OverallCode: ratelimit.RateLimitResponse_OK,
	}, nil
}

func NewTokenLimitServiceServer(tokenStore store.TokenStore) *TokenLimitServer {
	return &TokenLimitServer{tokenStore: tokenStore}
}
