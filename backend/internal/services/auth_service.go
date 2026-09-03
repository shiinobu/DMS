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

var ErrInvalidCredentials = errors.New("invalid username or password")

type AuthService struct {
	repository  *repositories.UserRepository
	jwtSecret   string
	jwtExpire   time.Duration
}

func NewAuthService(
	repository *repositories.UserRepository,
	jwtSecret string,
	jwtExpire time.Duration,
) *AuthService {
	return &AuthService{
		repository: repository,
		jwtSecret:  jwtSecret,
		jwtExpire:  jwtExpire,
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (string, *models.User, error) {

	username = strings.TrimSpace(username)

	user, err := s.repository.FindByUsername(
		ctx,
		username,
	)

	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return "", nil, ErrInvalidCredentials
		}

		return "", nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
		s.jwtSecret,
		s.jwtExpire,
	)

	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}