package seed

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

var roles = []struct {
	id          uuid.UUID
	name        string
	description string
}{
	{
		id:          uuid.MustParse("f8a5b4b0-4b1a-4b0e-8b0a-4b1a4b0e8b0a"),
		name:        "admin",
		description: "Super administrator with all permissions.",
	},
	{
		id:          uuid.MustParse("e7a4b3a9-3a09-4a9d-9a8c-3a094a9d9a8c"),
		name:        "engineer",
		description: "Engineer role with access to manage projects and environments.",
	},
}

func (s *Seeder) seedRoles(ctx context.Context) error {
	for _, r := range roles {
		var existingID uuid.UUID
		err := s.db.QueryRow(ctx, "SELECT id FROM roles WHERE name = $1", r.name).Scan(&existingID)
		if err == nil {
			log.Info().Str("name", r.name).Msg("role already exists, skipping")
			continue
		}

		if err != pgx.ErrNoRows {
			return err
		}

		_, err = s.db.Exec(ctx, `
			INSERT INTO roles (id, name, description)
			VALUES ($1, $2, $3)
		`, r.id, r.name, r.description)

		if err != nil {
			return err
		}

		log.Info().Str("name", r.name).Msg("seeded role")
	}
	return nil
}
