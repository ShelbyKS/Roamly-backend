package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"time"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type UserService struct {
	userStorage    storage.IUserStorage
	sessionStorage storage.ISessionStorage
}

func NewUserService(userStorage storage.IUserStorage, sessionStorage storage.ISessionStorage) service.IUserService {
	return &UserService{
		userStorage:    userStorage,
		sessionStorage: sessionStorage,
	}
}

func (s *UserService) Login(ctx context.Context, user model.User) (model.Session, error) {
	expectedUser, err := s.userStorage.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return model.Session{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	if !s.matchPasswords(expectedUser.Password, user.Password) {
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

func (s *UserService) Logout(ctx context.Context, session model.Session) error {
	err := s.sessionStorage.DeleteByToken(ctx, session.Token)
	if err != nil {
		return fmt.Errorf("failed to delete session in storage: %w", err)
	}

	return nil
}

func (s *UserService) Register(ctx context.Context, user model.User) (model.User, error) {
	_, err := s.userStorage.GetUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return model.User{}, fmt.Errorf("failed to check existing user: %w", err)
	}

	if !errors.Is(err, domain.ErrUserNotFound) {
		return model.User{}, fmt.Errorf("user already exists")
	}

	salt := make([]byte, 8)
	_, err = rand.Read(salt)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to generate salt: %w", err)
	}
	user.Password = s.hashPassword(salt, user.Password)

	err = s.userStorage.CreateUser(ctx, &user)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to create user in storage: %w", err)
	}

	return user, nil
}

func (s *UserService) hashPassword(salt, password []byte) []byte {
	hashedPassword := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)
	return append(salt, hashedPassword...)
}

func (s *UserService) matchPasswords(hashedPassword, plainPassword []byte) bool {
	salt := hashedPassword[:8]
	userPassHash := s.hashPassword(salt, plainPassword)

	return bytes.Equal(userPassHash, hashedPassword)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (model.User, error) {
	user, err := s.userStorage.GetUserByID(ctx, id)
	if err != nil {
		return model.User{}, fmt.Errorf("fail to get user from storage: %w", err)
	}

	return user, nil
}

func (service *UserService) UpdateUser(ctx context.Context, user model.User) error {
	err := service.userStorage.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("fail to update user from storage: %w", err)
	}

	return nil
}
