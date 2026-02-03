package seed

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func (s *Seeder) seedEnvironments(ctx context.Context) error {
	for _, p := range projects {
		var projectID uuid.UUID
		err := s.db.QueryRow(ctx, "SELECT id FROM projects WHERE slug = $1", p.slug).Scan(&projectID)
		if err != nil {
			log.Warn().Str("project_slug", p.slug).Err(err).Msg("could not find project to seed environment for")
			continue
		}

		envs := []struct {
			id   uuid.UUID
			name string
		}{
			{uuid.New(), "staging"},
			{uuid.New(), "production"},
		}

		for _, e := range envs {
			var existingID uuid.UUID
			queryErr := s.db.QueryRow(ctx, "SELECT id FROM environments WHERE project_id = $1 AND name = $2", projectID, e.name).Scan(&existingID)

			if queryErr == nil {
				log.Info().Str("project_slug", p.slug).Str("env_name", e.name).Msg("environment already exists, skipping")
				continue
			}

			if queryErr != pgx.ErrNoRows {
				return queryErr
			}

			_, execErr := s.db.Exec(ctx, `
				INSERT INTO environments (id, project_id, name, status, description)
				VALUES ($1, $2, $3, 'active', 'Seeded by Synthema seeder')
			`, e.id, projectID, e.name)

			if execErr != nil {
				return execErr
			}

			log.Info().Str("project_slug", p.slug).Str("env_name", e.name).Msg("seeded environment")
		}
	}
	return nil
}
