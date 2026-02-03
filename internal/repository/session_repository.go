package repository

import (
	"context"
	"database/sql"

	"synthema/internal/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	FindByID(ctx context.Context, id string) (*domain.Session, error)
	Delete(ctx context.Context, id string) error
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepository{db: db}
}

type sessionRepository struct {
	db *sql.DB
}

func (r *sessionRepository) Create(ctx context.Context, session *domain.Session) error {
	query := `INSERT INTO sessions (id, user_id, expires_at, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, session.ID, session.UserID, session.ExpiresAt, session.CreatedAt)
	return err
}

func (r *sessionRepository) FindByID(ctx context.Context, id string) (*domain.Session, error) {
	query := `SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var s domain.Session
	if err := row.Scan(&s.ID, &s.UserID, &s.ExpiresAt, &s.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}

func (r *sessionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
