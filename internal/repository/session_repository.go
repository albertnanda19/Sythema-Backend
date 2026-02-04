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
		INSERT INTO auth_sessions (id, user_id, expires_at)
		VALUES ($1, $2, $3)
		RETURNING created_at
	`
	return r.db.QueryRowContext(ctx, query, session.ID, session.UserID, session.ExpiresAt).Scan(&session.CreatedAt)
}

func (r *sessionRepository) FindActiveByID(ctx context.Context, id uuid.UUID) (*domain.AuthSession, error) {
	query := `
		SELECT id, user_id, expires_at, revoked_at, created_at
		FROM auth_sessions
		WHERE id = $1 AND revoked_at IS NULL AND expires_at > now()
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var (
		s         domain.AuthSession
		revokedAt sql.NullTime
	)
	if err := row.Scan(&s.ID, &s.UserID, &s.ExpiresAt, &revokedAt, &s.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if revokedAt.Valid {
		t := revokedAt.Time
		s.RevokedAt = &t
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
