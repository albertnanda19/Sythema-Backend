BEGIN;

ALTER TABLE auth_sessions
    DROP COLUMN IF EXISTS user_agent,
    DROP COLUMN IF EXISTS ip_address;

COMMIT;
