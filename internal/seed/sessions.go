package seed

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func (s *Seeder) seedSessions(ctx context.Context) error {
	var adminUserID uuid.UUID
	err := s.db.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", "admin@synthema.io").Scan(&adminUserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Warn().Msg("admin user not found, skipping session seeding")
			return nil
		}
		return err
	}

	var existingSessionID uuid.UUID
	err = s.db.QueryRow(ctx, "SELECT id FROM auth_sessions WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > now()", adminUserID).Scan(&existingSessionID)
	if err == nil {
		log.Info().Str("user_id", adminUserID.String()).Msg("active session for admin user already exists, skipping")
		return nil
	}

	if err != pgx.ErrNoRows {
		return err
	}

	sessionID := uuid.MustParse("a6a5b4a3-4a87-4a7c-9a6b-4a874a7c9a6b")
	expiresAt := time.Now().Add(24 * 30 * time.Hour) // 30 days

	_, err = s.db.Exec(ctx, `
		INSERT INTO auth_sessions (id, user_id, expires_at)
		VALUES ($1, $2, $3)
	`, sessionID, adminUserID, expiresAt)

	if err != nil {
		return err
	}

	log.Info().Str("user_id", adminUserID.String()).Msg("seeded session for admin user")
	return nil
}
