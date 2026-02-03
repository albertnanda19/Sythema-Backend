package seed

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

var projects = []struct {
	id   uuid.UUID
	slug string
	name string
}{
	{
		id:   uuid.MustParse("a9a7e72b-5f6a-4b0a-9b8e-3e4c5f6a7b8c"),
		slug: "payment-gateway",
		name: "Payment Gateway",
	},
	{
		id:   uuid.MustParse("b8b6e61a-4e59-4a99-8a7d-2d3b4e5f6a7b"),
		slug: "inventory-service",
		name: "Inventory Service",
	},
}

func (s *Seeder) seedProjects(ctx context.Context) error {
	for _, p := range projects {
		var existingID uuid.UUID
		err := s.db.QueryRow(ctx, "SELECT id FROM projects WHERE slug = $1", p.slug).Scan(&existingID)
		if err == nil {
			log.Info().Str("slug", p.slug).Msg("project already exists, skipping")
			continue
		}

		if err != pgx.ErrNoRows {
			return err
		}

		_, err = s.db.Exec(ctx, `
			INSERT INTO projects (id, slug, name, description)
			VALUES ($1, $2, $3, $4)
		`, p.id, p.slug, p.name, "Seeded by Synthema seeder")

		if err != nil {
			return err
		}

		log.Info().Str("slug", p.slug).Msg("seeded project")
	}
	return nil
}
