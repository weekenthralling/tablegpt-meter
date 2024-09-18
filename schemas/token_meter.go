package schemas

// UserToken represents the token information for a user.
type UserToken struct {
	UserID      string
	TotalTokens int32
	UsedTokens  int32
}
