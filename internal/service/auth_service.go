package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"synthema/internal/domain"
	"synthema/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// AuthService provides authentication services.
type AuthService interface {
	Authenticate(ctx context.Context, email, password string) (*domain.Session, error)
	Invalidate(ctx context.Context, sessionID string) error
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
func (s *authService) Authenticate(ctx context.Context, email, password string) (*domain.Session, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	session := &domain.Session{
		ID:        sessionID,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// Invalidate invalidates a session.
func (s *authService) Invalidate(ctx context.Context, sessionID string) error {
	return s.sessionRepo.Delete(ctx, sessionID)
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
