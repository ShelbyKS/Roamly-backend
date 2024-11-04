package service

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type AuthService struct {
	userStorage    storage.IUserStorage
	sessionStorage storage.ISessionStorage
}

func NewAuthService(userStorage storage.IUserStorage, sessionStorage storage.ISessionStorage) service.IAuthService {
	return &AuthService{
		userStorage:    userStorage,
		sessionStorage: sessionStorage,
	}
}

func (s *AuthService) Login(ctx context.Context, user model.User) (model.Session, error) {
	expectedUser, err := s.userStorage.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return model.Session{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	res, err := s.verifyPassword(expectedUser.Password, user.Password)
	if err != nil {
		return model.Session{}, fmt.Errorf("failed to verify password: %w", err)
	}
	if res == 0 {
		return model.Session{}, domain.ErrWrongCredentials
	}

	session := model.Session{
		Token:     uuid.NewString(),
		UserID:    expectedUser.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	err = s.sessionStorage.Add(ctx, session)
	if err != nil {
		return model.Session{}, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (s *AuthService) Logout(ctx context.Context, session model.Session) error {
	err := s.sessionStorage.DeleteByToken(ctx, session.Token)
	if err != nil {
		return fmt.Errorf("failed to delete session in storage: %w", err)
	}

	return nil
}

func (s *AuthService) Register(ctx context.Context, user model.User) (model.User, error) {
	_, err := s.userStorage.GetUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return model.User{}, fmt.Errorf("failed to check existing user: %w", err)
	}

	if !errors.Is(err, domain.ErrUserNotFound) {
		return model.User{}, domain.ErrUserAlreadyExists
	}

	salt, err := s.generateSalt()
	if err != nil {
		return model.User{}, fmt.Errorf("failed to generate salt: %w", err)
	}

	user.Password = s.hashPassword(user.Password, salt)

	err = s.userStorage.CreateUser(ctx, &user)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to create user in storage: %w", err)
	}

	return user, nil
}

func (s *AuthService) generateSalt() ([]byte, error) {
	salt := make([]byte, 16) // Например, 16 байт
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func (s *AuthService) hashPassword(password string, salt []byte) string {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32) // Параметры Argon2id
	return fmt.Sprintf("%s:%s", base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash))
}

func (s *AuthService) verifyPassword(storedHash, password string) (int, error) {
	parts := strings.Split(storedHash, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid stored hash format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, err
	}
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return subtle.ConstantTimeCompare(hash, expectedHash), nil
}
