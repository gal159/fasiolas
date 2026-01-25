# 🔧 FIXES APPLIED - Column Name Mismatch

## Issues Found and Fixed

### Issue 1: game_actions table schema mismatch ✅ FIXED
**Error:** `pq: column "player_id" does not exist`

**Root Cause:** 
- Migration `000004_create_game_actions_table.up.sql` defines: `user_id` and `action_data`, `phase`, `timestamp`
- Repository `game_action_repository.go` was querying: `player_id` (wrong column name)
- init-db.sh had outdated schema definition

**Solutions Applied:**
1. ✅ Updated `init-db.sh` - game_actions table now includes `user_id`, `action_data`, `phase`, `timestamp`
2. ✅ Fixed `game_action_repository.go`:
   - Line 21: Changed `player_id` → `user_id` in INSERT query
   - Line 44: Changed `player_id` → `user_id` in SELECT for GetByGame
   - Line 88: Changed `player_id` → `user_id` in SELECT for GetByUser

### Changes Made

**File: `internal/repository/game_action_repository.go`**
- Create method: `INSERT INTO game_actions (game_id, user_id, ...)` ✅
- GetByGame method: `SELECT id, game_id, user_id, ...` ✅  
- GetByUser method: `WHERE user_id = $1` ✅

**File: `init-db.sh`**
- Created game_actions with correct columns: `user_id`, `action_type`, `action_data`, `phase`, `timestamp`
- Added constraint check for valid action types
- Added proper indexes

## Result
The game should now properly:
- ✅ Create games
- ✅ Load game state
- ✅ Track game actions
- ✅ Fetch action history

## How to Test
1. Refresh browser or restart docker: `docker-compose restart app`
2. Create a new game
3. The 400 error should be gone
4. Game board should load with players and deck visible
