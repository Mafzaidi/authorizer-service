package repository

import "context"

type AuthRepository interface {
	StoreRefreshToken(ctx context.Context, userID, token string) error
	GetRefreshToken(ctx context.Context, userID string) (string, error)
	MapRefreshToUser(ctx context.Context, refreshToken, userID string) error
	GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (string, error)
	DeleteRefreshToken(ctx context.Context, userID string) error
}
