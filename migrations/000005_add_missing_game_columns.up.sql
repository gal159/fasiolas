-- Add missing columns to games table if they don't exist
ALTER TABLE games
ADD COLUMN IF NOT EXISTS started_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS finished_at TIMESTAMP WITH TIME ZONE;
