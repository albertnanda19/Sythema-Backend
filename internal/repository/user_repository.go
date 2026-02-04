package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"synthema/internal/domain"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	ListRolesByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Role, error)
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	row := r.db.QueryRowContext(ctx, query, email)

	var (
		u         domain.User
		deletedAt sql.NullTime
	)
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		u.DeletedAt = &t
	}

	return &u, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var (
		u         domain.User
		deletedAt sql.NullTime
	)
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		u.DeletedAt = &t
	}

	return &u, nil
}

func (r *userRepository) ListRolesByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
		ORDER BY r.name
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]domain.Role, 0)
	for rows.Next() {
		var (
			r           domain.Role
			description sql.NullString
			createdAt   time.Time
			updatedAt   time.Time
		)
		if err := rows.Scan(&r.ID, &r.Name, &description, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		r.CreatedAt = createdAt
		r.UpdatedAt = updatedAt
		if description.Valid {
			d := description.String
			r.Description = &d
		}
		roles = append(roles, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}
