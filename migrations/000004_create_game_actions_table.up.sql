-- Create game_actions table (for audit log and game history)
CREATE TABLE IF NOT EXISTS game_actions (
    id SERIAL PRIMARY KEY,
    game_id INTEGER NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    action_type VARCHAR(50) NOT NULL,
    action_data JSONB,
    phase INTEGER,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_action_type CHECK (action_type IN (
        'join', 'leave', 'start_game', 'draw_card', 'place_card',
        'call_cheat', 'receive_penalty', 'phase_change',
        'take_table_card', 'skip_turn', 'win', 'lose'
    ))
);

CREATE INDEX idx_game_actions_game ON game_actions(game_id);
CREATE INDEX idx_game_actions_user ON game_actions(user_id);
CREATE INDEX idx_game_actions_timestamp ON game_actions(timestamp);
CREATE INDEX idx_game_actions_type ON game_actions(action_type);

