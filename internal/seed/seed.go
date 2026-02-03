package seed

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Seeder struct {
	db *pgxpool.Pool
}

func NewSeeder(db *pgxpool.Pool) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) Run(ctx context.Context) error {
	seedFuncs := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"projects", s.seedProjects},
		{"environments", s.seedEnvironments},
		{"roles", s.seedRoles},
		{"users and user_roles", s.seedUsersAndRoles},
		{"sessions", s.seedSessions},
	}

	for _, f := range seedFuncs {
		log.Info().Msgf("starting to seed %s...", f.name)
		if err := f.fn(ctx); err != nil {
			log.Error().Err(err).Msgf("failed to seed %s", f.name)
			return err
		}
		log.Info().Msgf("finished seeding %s.", f.name)
	}

	return nil
}
