package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"synthema/internal/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, session *domain.AuthSession) error
	FindActiveByID(ctx context.Context, id uuid.UUID) (*domain.AuthSession, error)
	Revoke(ctx context.Context, id uuid.UUID, revokedAt time.Time) error
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepository{db: db}
}

type sessionRepository struct {
	db *sql.DB
}

func (r *sessionRepository) Create(ctx context.Context, session *domain.AuthSession) error {
	query := `
		INSERT INTO auth_sessions (id, user_id, expires_at, user_agent, ip_address)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`
	return r.db.QueryRowContext(ctx, query, session.ID, session.UserID, session.ExpiresAt, session.UserAgent, session.IPAddress).Scan(&session.CreatedAt)
}

func (r *sessionRepository) FindActiveByID(ctx context.Context, id uuid.UUID) (*domain.AuthSession, error) {
	query := `
		SELECT id, user_id, expires_at, revoked_at, created_at, user_agent, ip_address
		FROM auth_sessions
		WHERE id = $1 AND revoked_at IS NULL AND expires_at > now()
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var (
		s         domain.AuthSession
		revokedAt sql.NullTime
		userAgent sql.NullString
		ipAddress sql.NullString
	)
	if err := row.Scan(&s.ID, &s.UserID, &s.ExpiresAt, &revokedAt, &s.CreatedAt, &userAgent, &ipAddress); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if revokedAt.Valid {
		t := revokedAt.Time
		s.RevokedAt = &t
	}
	if userAgent.Valid {
		ua := userAgent.String
		s.UserAgent = &ua
	}
	if ipAddress.Valid {
		ip := ipAddress.String
		s.IPAddress = &ip
	}

	return &s, nil
}

func (r *sessionRepository) Revoke(ctx context.Context, id uuid.UUID, revokedAt time.Time) error {
	query := `
		UPDATE auth_sessions
		SET revoked_at = $2
		WHERE id = $1 AND revoked_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, id, revokedAt)
	return err
}
