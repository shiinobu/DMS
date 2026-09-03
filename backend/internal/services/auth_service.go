package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"dms/backend/internal/auth"
	"dms/backend/internal/models"
	"dms/backend/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New(
	"invalid username or password",
)

type AuthService struct {
	repository    *repositories.UserRepository
	jwtSecret     string
	jwtExpiration time.Duration
}

func NewAuthService(
	repository *repositories.UserRepository,
	jwtSecret string,
	jwtExpiration time.Duration,
) *AuthService {
	return &AuthService{
		repository:    repository,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

type LoginResult struct {
	Token string
	User  *models.User
}

func (s *AuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (*LoginResult, error) {
	username = strings.TrimSpace(username)

	user, err := s.repository.FindByUsername(
		ctx,
		username,
	)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
		s.jwtSecret,
		s.jwtExpiration,
	)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Token: token,
		User:  user,
	}, nil
}
