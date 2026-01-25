-- Fasiolas Card Game - All Database Migrations
-- This file contains all 4 database migrations in order

-- ===========================================
-- Migration 1: Create users table
-- ===========================================
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

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_oauth ON users(oauth_provider, oauth_id);
CREATE INDEX idx_users_role ON users(role);

-- ===========================================
-- Migration 2: Create games table
-- ===========================================
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
    winner_id INTEGER,
    CONSTRAINT fk_games_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_games_winner FOREIGN KEY (winner_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_games_room_code ON games(room_code);
CREATE INDEX idx_games_state ON games(state);
CREATE INDEX idx_games_created_by ON games(created_by);

-- ===========================================
-- Migration 3: Create game_players table
-- ===========================================
CREATE TABLE IF NOT EXISTS game_players (
    id SERIAL PRIMARY KEY,
    game_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    position INTEGER NOT NULL,
    cards JSONB DEFAULT '[]',
    top_card VARCHAR(10),
    card_count INTEGER DEFAULT 0,
    status VARCHAR(50) DEFAULT 'waiting',
    can_call_cheat BOOLEAN DEFAULT TRUE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_game_players_game FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
    CONSTRAINT fk_game_players_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_game_position UNIQUE(game_id, position)
);

CREATE INDEX idx_game_players_game_id ON game_players(game_id);
CREATE INDEX idx_game_players_user_id ON game_players(user_id);

-- ===========================================
-- Migration 4: Create game_actions table
-- ===========================================
CREATE TABLE IF NOT EXISTS game_actions (
    id SERIAL PRIMARY KEY,
    game_id INTEGER NOT NULL,
    player_id INTEGER NOT NULL,
    action_type VARCHAR(50) NOT NULL,
    card_placed VARCHAR(10),
    target_player_position INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_game_actions_game FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
    CONSTRAINT fk_game_actions_player FOREIGN KEY (player_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_game_actions_game_id ON game_actions(game_id);
CREATE INDEX idx_game_actions_player_id ON game_actions(player_id);
CREATE INDEX idx_game_actions_created_at ON game_actions(created_at);

-- ===========================================
-- Verification query:
-- SELECT * FROM information_schema.tables WHERE table_schema='public';
-- ===========================================
