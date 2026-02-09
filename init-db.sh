#!/bin/bash
# Database initialization script - runs migrations on startup

set -e

echo "=== Starting Database Initialization ==="
echo "Waiting for PostgreSQL to be ready..."

# Wait for PostgreSQL to be ready
for i in {1..30}; do
    if PGPASSWORD=123456 pg_isready -h postgres -U postgres > /dev/null 2>&1; then
        echo "✓ PostgreSQL is ready"
        break
    fi
    echo "  Attempt $i/30: Waiting for PostgreSQL..."
    sleep 1
done

# Create database if it doesn't exist
echo "Creating database if not exists..."
PGPASSWORD=123456 psql -h postgres -U postgres -c "CREATE DATABASE IF NOT EXISTS fasiolas_game;" || true

# Run migrations
echo "Applying migrations..."

echo "  [1/4] Creating users table..."
PGPASSWORD=123456 psql -h postgres -U postgres -d fasiolas_game << 'EOF'
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    oauth_provider VARCHAR(50) NOT NULL,
    oauth_id VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'player',
    avatar_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP WITH TIME ZONE,
    CONSTRAINT unique_oauth UNIQUE(oauth_provider, oauth_id),
    CONSTRAINT check_role CHECK (role IN ('admin', 'player', 'spectator'))
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_oauth ON users(oauth_provider, oauth_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
EOF

echo "  [2/4] Creating games table..."
PGPASSWORD=123456 psql -h postgres -U postgres -d fasiolas_game << 'EOF'
CREATE TABLE IF NOT EXISTS games (
    id SERIAL PRIMARY KEY,
    room_code VARCHAR(10) UNIQUE NOT NULL,
    state VARCHAR(50) NOT NULL DEFAULT 'waiting',
    phase INTEGER DEFAULT 1,
    current_player_position INTEGER DEFAULT 0,
    trump_suit VARCHAR(20),
    deck_cards JSONB DEFAULT '[]',
    table_cards JSONB DEFAULT '[]',
    max_players INTEGER DEFAULT 2,
    created_by INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP WITH TIME ZONE,
    finished_at TIMESTAMP WITH TIME ZONE,
    winner_id INTEGER,
    CONSTRAINT fk_games_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_games_winner FOREIGN KEY (winner_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Add missing columns safely if table already existed
ALTER TABLE games ADD COLUMN IF NOT EXISTS started_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE games ADD COLUMN IF NOT EXISTS finished_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX IF NOT EXISTS idx_games_room_code ON games(room_code);
CREATE INDEX IF NOT EXISTS idx_games_state ON games(state);
CREATE INDEX IF NOT EXISTS idx_games_created_by ON games(created_by);
EOF

echo "  [3/4] Creating game_players table..."
PGPASSWORD=123456 psql -h postgres -U postgres -d fasiolas_game << 'EOF'
CREATE TABLE IF NOT EXISTS game_players (
    id SERIAL PRIMARY KEY,
    game_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    position INTEGER NOT NULL,
    cards JSONB DEFAULT '[]',
    top_card JSONB,
    card_count INTEGER DEFAULT 0,
    status VARCHAR(50) DEFAULT 'waiting',
    can_call_cheat BOOLEAN DEFAULT TRUE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_game_players_game FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
    CONSTRAINT fk_game_players_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_game_position UNIQUE(game_id, position)
);

CREATE INDEX IF NOT EXISTS idx_game_players_game_id ON game_players(game_id);
CREATE INDEX IF NOT EXISTS idx_game_players_user_id ON game_players(user_id);
EOF

echo "  [4/4] Creating game_actions table..."
PGPASSWORD=123456 psql -h postgres -U postgres -d fasiolas_game << 'EOF'
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

CREATE INDEX IF NOT EXISTS idx_game_actions_game ON game_actions(game_id);
CREATE INDEX IF NOT EXISTS idx_game_actions_user ON game_actions(user_id);
CREATE INDEX IF NOT EXISTS idx_game_actions_timestamp ON game_actions(timestamp);
CREATE INDEX IF NOT EXISTS idx_game_actions_type ON game_actions(action_type);
EOF

echo "✓ All migrations applied successfully"
echo "=== Database Initialization Complete ==="
