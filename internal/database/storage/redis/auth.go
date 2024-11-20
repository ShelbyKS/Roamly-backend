package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"

	"github.com/redis/go-redis/v9"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type SessionStorage struct {
	client *redis.Client
}

func NewSessionStorage(client *redis.Client) storage.ISessionStorage {
	return &SessionStorage{
		client: client,
	}
}

func (s *SessionStorage) Add(ctx context.Context, session model.Session) error {
	if session.Token == "" {
		return fmt.Errorf("session token is empty")
	}
	duration := session.ExpiresAt.Sub(time.Now())
	err := s.client.Set(ctx, session.Token, session.UserID, duration).Err()
	if err != nil {
		return fmt.Errorf("failed to add session: %w", err)
	}

	userKey := fmt.Sprintf("user:%d", session.UserID)

	err = s.client.LPush(ctx, userKey, session.Token).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to user list: %w", err)
	}

	err = s.client.Expire(ctx, userKey, duration).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration for user list: %w", err)
	}

	return nil
}

func (s *SessionStorage) GetTokensByUserID(ctx context.Context, userID int) ([]string, error) {
	userKey := fmt.Sprintf("user:%d", userID)

	tokens, err := s.client.LRange(ctx, userKey, 0, -1).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get tokens for user %d: %w", userID, err)
	}

	return tokens, nil
}


func (s *SessionStorage) DeleteByToken(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("session token is empty")
	}
	err := s.client.Del(ctx, token).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

func (s *SessionStorage) SessionExists(ctx context.Context, token string) (model.Session, error) {
	if token == "" {
		return model.Session{}, fmt.Errorf("session token is empty")
	}

	userIDStr, err := s.client.Get(ctx, token).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Session{}, domain.ErrSessionNotFound
		}
		return model.Session{}, fmt.Errorf("failed to get session: %w", err)
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return model.Session{}, fmt.Errorf("failed to convert userID to int: %w", err)
	}

	return model.Session{
		Token:  token,
		UserID: userID,
	}, nil
}
