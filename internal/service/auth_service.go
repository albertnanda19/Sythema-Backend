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

type SessionMeta struct {
	UserAgent string
	IPAddress string
}

// AuthService provides authentication services.
type AuthService interface {
	Authenticate(ctx context.Context, email, password string, meta SessionMeta) (*domain.User, *domain.AuthSession, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
}

// NewAuthService creates a new auth service.
func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, sessionTTL time.Duration) AuthService {
	if sessionTTL <= 0 {
		sessionTTL = 7 * 24 * time.Hour
	}
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		sessionTTL:  sessionTTL,
	}
}

type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	sessionTTL  time.Duration
}

// Authenticate authenticates a user and creates a session.
func (s *authService) Authenticate(ctx context.Context, email, password string, meta SessionMeta) (*domain.User, *domain.AuthSession, error) {
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

	roles, err := s.userRepo.ListRolesByUserID(ctx, user.ID)
	if err != nil {
		return nil, nil, err
	}
	user.Roles = roles

	now := time.Now()

	session := &domain.AuthSession{
		ID:        uuid.New(),
		UserID:    user.ID,
		ExpiresAt: now.Add(s.sessionTTL),
	}
	if meta.UserAgent != "" {
		ua := meta.UserAgent
		session.UserAgent = &ua
	}
	if meta.IPAddress != "" {
		ip := meta.IPAddress
		session.IPAddress = &ip
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
