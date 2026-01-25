-- Remove added columns (rollback)
ALTER TABLE games
DROP COLUMN IF EXISTS started_at,
DROP COLUMN IF EXISTS finished_at;
