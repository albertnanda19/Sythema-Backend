package repositories

import (
	"database/sql"

	"synthema/internal/repository"
)

type SessionRepository = repository.SessionRepository

func NewSessionRepository(db *sql.DB) SessionRepository {
	return repository.NewSessionRepository(db)
}
