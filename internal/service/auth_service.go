package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"synthema/internal/domain"
	"synthema/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// AuthService provides authentication services.
type AuthService interface {
	Authenticate(ctx context.Context, email, password string) (*domain.User, *domain.AuthSession, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
}

// NewAuthService creates a new auth service.
func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

// Authenticate authenticates a user and creates a session.
func (s *authService) Authenticate(ctx context.Context, email, password string) (*domain.User, *domain.AuthSession, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, ErrInvalidCredentials
	}
	if !user.IsActive {
		return nil, nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	session := &domain.AuthSession{
		ID:        uuid.New(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

// RevokeSession revokes a session.
func (s *authService) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.sessionRepo.Revoke(ctx, sessionID, time.Now())
}
