-- Create game_players table
CREATE TABLE IF NOT EXISTS game_players (
    id SERIAL PRIMARY KEY,
    game_id INTEGER NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    cards JSONB DEFAULT '[]'::jsonb,
    top_card JSONB,
    card_count INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    can_call_cheat BOOLEAN DEFAULT true,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_game_user UNIQUE(game_id, user_id),
    CONSTRAINT unique_game_position UNIQUE(game_id, position),
    CONSTRAINT check_status CHECK (status IN ('active', 'eliminated', 'winner', 'left'))
);

CREATE INDEX idx_game_players_game ON game_players(game_id);
CREATE INDEX idx_game_players_user ON game_players(user_id);
CREATE INDEX idx_game_players_position ON game_players(game_id, position);

