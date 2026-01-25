-- Create games table
CREATE TABLE IF NOT EXISTS games (
    id SERIAL PRIMARY KEY,
    room_code VARCHAR(10) UNIQUE NOT NULL,
    state VARCHAR(20) NOT NULL DEFAULT 'waiting',
    phase INTEGER DEFAULT 1,
    current_player_position INTEGER,
    trump_suit VARCHAR(10),
    deck_cards JSONB DEFAULT '[]'::jsonb,
    table_cards JSONB DEFAULT '[]'::jsonb,
    created_by INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP WITH TIME ZONE,
    finished_at TIMESTAMP WITH TIME ZONE,
    winner_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    max_players INTEGER DEFAULT 8 CHECK (max_players >= 2 AND max_players <= 8),
    CONSTRAINT fk_games_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT check_state CHECK (state IN ('waiting', 'phase1', 'phase2', 'finished', 'cancelled')),
    CONSTRAINT check_phase CHECK (phase IN (1, 2))
);

CREATE INDEX idx_games_state ON games(state);
CREATE INDEX idx_games_room_code ON games(room_code);
CREATE INDEX idx_games_created_by ON games(created_by);

