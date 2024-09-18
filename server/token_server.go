package server

import (
	"context"

	token "github.com/tablegpt_meter/proto/token"
	"github.com/tablegpt_meter/store"
)

// TokenServiceServer implements the TokenService gRPC service.
type TokenServiceServer struct {
	token.TokenServiceServer
	tokenStore store.TokenStore
}

// NewTokenServiceServer creates a new TokenServiceServer instance.
func NewTokenServiceServer(store store.TokenStore) *TokenServiceServer {
	return &TokenServiceServer{
		tokenStore: store,
	}
}

// RecordTokenUsage creates a record of token usage for a specific user.
func (s *TokenServiceServer) RecordTokenUsage(ctx context.Context, req *token.RecordTokenUsageRequest) (*token.TokenOperationResponse, error) {
	err := s.tokenStore.SaveTokenUsageRecord(ctx, req.UserId, req.Domain, int32(req.TokensUsed))
	if err != nil {
		return &token.TokenOperationResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	return &token.TokenOperationResponse{
		Success: true,
		Message: "Token usage recorded successfully.",
	}, nil
}

// UpdateUserTotalTokens updates the total number of tokens for a specific user.
func (s *TokenServiceServer) UpdateUserTotalTokens(ctx context.Context, req *token.UpdateUserTotalTokensRequest) (*token.TokenOperationResponse, error) {
	err := s.tokenStore.SaveUserToken(ctx, req.UserId, req.Domain, int32(req.NewTotalTokens))
	if err != nil {
		return &token.TokenOperationResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	return &token.TokenOperationResponse{
		Success: true,
		Message: "Total tokens updated successfully.",
	}, nil
}
