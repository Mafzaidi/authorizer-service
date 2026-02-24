package repository

import (
	"context"
	"time"

	"github.com/mafzaidi/authorizer/internal/domain/repository"
	"github.com/redis/go-redis/v9"
)

type authRepository struct {
	redis *redis.Client
}

func NewAuthRepository(redis *redis.Client) repository.AuthRepository {
	return &authRepository{
		redis: redis,
	}
}

func (r *authRepository) StoreRefreshToken(ctx context.Context, userID, token string) error {
	return r.redis.Set(ctx, "refresh:"+userID, token, 7*24*time.Hour).Err()
}

func (r *authRepository) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	return r.redis.Get(ctx, "refresh:"+userID).Result()
}

func (r *authRepository) MapRefreshToUser(ctx context.Context, refreshToken, userID string) error {
	return r.redis.Set(ctx, "rt:"+refreshToken, userID, 7*24*time.Hour).Err()
}

func (r *authRepository) GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	return r.redis.Get(ctx, "rt:"+refreshToken).Result()
}

func (r *authRepository) DeleteRefreshToken(ctx context.Context, userID string) error {
	return r.redis.Del(ctx, "refresh:"+userID).Err()
}
