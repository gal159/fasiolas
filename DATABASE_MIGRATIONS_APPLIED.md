# ✅ DATABASE MIGRATIONS APPLIED

## Status: Complete ✅

All database migrations have been successfully applied to the PostgreSQL database.

### Tables Created

1. **users** - User accounts and OAuth data
   - id, email, username, oauth_provider, oauth_id
   - role (admin, player, spectator)
   - avatar_url, timestamps
   - Indexes: email, oauth, role

2. **games** - Game sessions
   - id, room_code, state (waiting/phase1/phase2/finished)
   - phase, current_player_position, trump_suit
   - deck_cards, table_cards, max_players
   - created_by, timestamps, winner_id

3. **game_players** - Players in games
   - id, game_id, user_id, position
   - cards, top_card, card_count, status
   - can_call_cheat, joined_at
   - Indexes: game_id, user_id

4. **game_actions** - Action history
   - id, game_id, user_id, action_type
   - action_data (JSON), phase, timestamp
   - Indexes: game_id, user_id, action_type, timestamp

### How to Apply Migrations

The migrations have been manually applied using the following commands:

```bash
# From the cardGame directory, run each migration in order:

# 1. Create users table
docker compose exec postgres psql -U postgres -d fasiolas_game < migrations/000001_create_users_table.up.sql

# 2. Create games table
docker compose exec postgres psql -U postgres -d fasiolas_game < migrations/000002_create_games_table.up.sql

# 3. Create game_players table
docker compose exec postgres psql -U postgres -d fasiolas_game < migrations/000003_create_game_players_table.up.sql

# 4. Create game_actions table
docker compose exec postgres psql -U postgres -d fasiolas_game < migrations/000004_create_game_actions_table.up.sql
```

Or using PowerShell (on Windows):

```powershell
# In the cardGame directory:
Get-Content migrations/000001_create_users_table.up.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game
Get-Content migrations/000002_create_games_table.up.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game
Get-Content migrations/000003_create_game_players_table.up.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game
Get-Content migrations/000004_create_game_actions_table.up.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game
```

### Verify Migrations

To verify all tables exist:

```bash
docker compose exec postgres psql -U postgres -d fasiolas_game -c "\dt"
```

Expected output:
```
            List of relations
 Schema      Name      Type    Owner   
--------+--------------+-------+----------
 public  game_actions  table  postgres
 public  game_players  table  postgres
 public  games         table  postgres
 public  users         table  postgres
(4 rows)
```

### Rolling Back

If you need to rollback migrations (use .down.sql files):

```bash
# In reverse order (4, 3, 2, 1):
Get-Content migrations/000004_create_game_actions_table.down.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game
Get-Content migrations/000003_create_game_players_table.down.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game
Get-Content migrations/000002_create_games_table.down.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game
Get-Content migrations/000001_create_users_table.down.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game
```

### Database Refresh

To completely reset the database (delete everything and recreate):

```bash
# 1. Delete the volume
docker compose down -v

# 2. Start fresh
docker compose up -d

# 3. Apply migrations again
Get-Content migrations/000001_create_users_table.up.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game
Get-Content migrations/000002_create_games_table.up.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game
Get-Content migrations/000003_create_game_players_table.up.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game
Get-Content migrations/000004_create_game_actions_table.up.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game
```

---

## What This Fixes

The error you were getting:
```
{"error":"failed to get user: pq: relation \"users\" does not exist"}
```

This error occurred because:
1. The database existed but had no tables
2. The migrations weren't applied automatically
3. Backend tried to query the `users` table which didn't exist

**Now that migrations are applied**, the backend can:
- ✅ Create user accounts
- ✅ Store OAuth data
- ✅ Create games
- ✅ Manage players
- ✅ Track game actions

---

## Next Steps

1. **Refresh the browser**: http://localhost:3000
2. **Login again**: Try Google login again
3. **Create a game**: You should now be able to create games
4. **Play**: Everything should work now!

---

## Status

✅ All 4 migrations applied
✅ All 4 tables created
✅ Backend restarted
✅ Ready to use

**You can now login and play! 🎮**

