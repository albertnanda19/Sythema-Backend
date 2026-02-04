package repositories

import (
	"database/sql"

	"synthema/internal/repository"
)

type UserRepository = repository.UserRepository

func NewUserRepository(db *sql.DB) UserRepository {
	return repository.NewUserRepository(db)
}
