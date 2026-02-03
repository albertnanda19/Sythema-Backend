package seed

import (
	"context"
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func (s *Seeder) seedUsersAndRoles(ctx context.Context) error {
	adminPassword := os.Getenv("SYNTHEMA_SEED_ADMIN_PASSWORD")
	if adminPassword == "" {
		return errors.New("environment variable SYNTHEMA_SEED_ADMIN_PASSWORD must be set for seeding")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []struct {
		id    uuid.UUID
		email string
		role  string
	}{
		{
			id:    uuid.MustParse("c7a6b5a4-5a98-4a8d-8a7c-5a984a8d8a7c"),
			email: "admin@synthema.io",
			role:  "admin",
		},
		{
			id:    uuid.MustParse("d6a5b4a3-4a87-4a7c-9a6b-4a874a7c9a6b"),
			email: "engineer@synthema.io",
			role:  "engineer",
		},
	}

	for _, u := range users {
		// Seed user
		var existingUserID uuid.UUID
		err := s.db.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", u.email).Scan(&existingUserID)
		if err == nil {
			log.Info().Str("email", u.email).Msg("user already exists, skipping user creation")
		} else if err == pgx.ErrNoRows {
			_, err = s.db.Exec(ctx, `
				INSERT INTO users (id, email, password_hash)
				VALUES ($1, $2, $3)
			`, u.id, u.email, string(hashedPassword))
			if err != nil {
				return err
			}
			log.Info().Str("email", u.email).Msg("seeded user")
			existingUserID = u.id
		} else {
			return err
		}

		// Seed user_role relationship
		var roleID uuid.UUID
		err = s.db.QueryRow(ctx, "SELECT id FROM roles WHERE name = $1", u.role).Scan(&roleID)
		if err != nil {
			log.Warn().Str("role", u.role).Msg("role not found, cannot assign to user")
			continue
		}

		var existingUserRoleID uuid.UUID
		err = s.db.QueryRow(ctx, "SELECT id FROM user_roles WHERE user_id = $1 AND role_id = $2", existingUserID, roleID).Scan(&existingUserRoleID)
		if err == nil {
			log.Info().Str("email", u.email).Str("role", u.role).Msg("user_role relationship already exists, skipping")
			continue
		}

		if err != pgx.ErrNoRows {
			return err
		}

		_, err = s.db.Exec(ctx, `
			INSERT INTO user_roles (id, user_id, role_id)
			VALUES ($1, $2, $3)
		`, uuid.New(), existingUserID, roleID)

		if err != nil {
			return err
		}
		log.Info().Str("email", u.email).Str("role", u.role).Msg("assigned role to user")
	}

	return nil
}
