# 🔧 Join Game Fix - "Already in Game" Error

## 📋 Problema

**Error:** `"Failed to join game: already in game"`

**Kada:** Kai vartotojas bando prisijungti prie žaidimo, net jei jis nėra tame žaidime.

---

## 🔍 Root Cause

Sistema tikrina ar vartotojas jau yra **tame konkrečiame žaidime**, bet **netikrina** ar vartotojas jau yra **kitame aktyviam žaidime**.

Tai reiškia:
- ✅ Vartotojas gali sukurti naują žaidimą (automatiškai pridedamas kaip žaidėjas)
- ❌ Vartotojas negali prisijungti prie kito žaidimo (jau yra pirmame žaidime)

**Pavyzdys:**
1. Vartotojas sukuria žaidimą A → pridedamas kaip žaidėjas
2. Vartotojas bando prisijungti prie žaidimo B → **ERROR**: "already in game"

---

## ✅ Sprendimas

### 1. Pridėtas naujas repository metodas: `GetUserActiveGame`

**Failas:** `internal/repository/game_player_repository.go`

```go
// GetUserActiveGame retrieves the active game for a user (if any)
// Returns the game player record if user is in an active game (waiting or playing)
func (r *GamePlayerRepository) GetUserActiveGame(userID int) (*models.GamePlayer, error) {
	gp := &models.GamePlayer{}
	var cardsJSON, topCardJSON []byte

	query := `
		SELECT gp.id, gp.game_id, gp.user_id, gp.position, gp.cards, gp.top_card, 
		       gp.card_count, gp.status, gp.can_call_cheat, gp.joined_at
		FROM game_players gp
		JOIN games g ON gp.game_id = g.id
		WHERE gp.user_id = $1 
		  AND g.state IN ('waiting', 'playing')
		  AND gp.status = 'active'
		LIMIT 1
	`
	// ... implementation
}
```

**Kas daro:**
- Ieško ar vartotojas yra **bet kuriame** žaidime
- Tikrina tik **aktyvius** žaidimus (state = 'waiting' arba 'playing')
- Tikrina tik **aktyvius** žaidėjus (status = 'active')
- Grąžina `nil` jei vartotojas nėra jokiame žaidime

---

### 2. Atnaujintas `JoinGame` metodas

**Failas:** `internal/service/game_service.go`

**Prieš:**
```go
// Check if user already in game
existing, err := s.playerRepo.GetByGameAndUser(g.ID, userID)
if err != nil {
	return nil, fmt.Errorf("failed to check existing player: %w", err)
}
if existing != nil {
	return nil, errors.New("already in game")
}
```

**Po:**
```go
// Check if user already in THIS specific game
existing, err := s.playerRepo.GetByGameAndUser(g.ID, userID)
if err != nil {
	return nil, fmt.Errorf("failed to check existing player: %w", err)
}
if existing != nil {
	return nil, errors.New("already in this game")
}

// Check if user is already in ANY active game
activeGame, err := s.playerRepo.GetUserActiveGame(userID)
if err != nil {
	return nil, fmt.Errorf("failed to check user active games: %w", err)
}
if activeGame != nil && activeGame.GameID != g.ID {
	return nil, errors.New("already in another active game")
}
```

**Kas pasikeitė:**
1. **Tikslesnis error message:** "already in this game" vs "already in another active game"
2. **Dvigubas patikrinimas:**
   - Pirma tikrina ar jau yra **tame pačiame** žaidime
   - Paskui tikrina ar yra **kitame** aktyviam žaidime
3. **Leidžia re-join:** Jei activeGame.GameID == g.ID (tas pats žaidimas), leidžia

---

## 📊 Validation Logic

```
Vartotojas bando prisijungti prie žaidimo
↓
Tikrinti:
├─ 1. Ar žaidimas egzistuoja? ✅
├─ 2. Ar žaidimas dar neprasidėjęs? (state = 'waiting') ✅
├─ 3. Ar vartotojas jau yra ŠIAME žaidime? ❌
├─ 4. Ar vartotojas yra KITAME aktyviam žaidime? ❌
├─ 5. Ar žaidimas pilnas? ❌
└─ ✅ Pridėti vartotoją prie žaidimo
```

---

## 🧪 Test Scenarios

### Scenario 1: Vartotojas sukuria žaidimą ir bando prisijungti prie to paties
**Rezultatas:** ❌ "already in this game"
**Teisingai:** Vartotojas jau yra tame žaidime, nereikia join

---

### Scenario 2: Vartotojas sukuria žaidimą A ir bando prisijungti prie žaidimo B
**Rezultatas:** ❌ "already in another active game"
**Teisingai:** Vartotojas negali būti dviejuose žaidimuose vienu metu

---

### Scenario 3: Vartotojas prisijungė prie žaidimo A, žaidimas baigėsi, dabar join prie B
**Rezultatas:** ✅ Success
**Teisingai:** Vartotojas dabar nėra aktyviam žaidime (state = 'finished')

---

### Scenario 4: Vartotojas dar nėra jokiame žaidime, prisijungia prie A
**Rezultatas:** ✅ Success
**Teisingai:** Pirmas join, viskas gerai

---

## 🔄 User Flow (Pataisytas)

### Normalus Flow:
```
1. Vartotojas eina į /dashboard
2. Mato laukiančius žaidimus (state = 'waiting')
3. Pasirenka žaidimą
4. Paspauskia "Join Game"
5. Sistema tikrina:
   ✅ Žaidimas egzistuoja
   ✅ Žaidimas dar laukia žaidėjų
   ✅ Vartotojas dar nėra tame žaidime
   ✅ Vartotojas nėra kitame aktyviam žaidime
   ✅ Žaidimas nepilnas
6. ✅ Vartotojas pridedamas prie žaidimo
7. ✅ Vartotojas nukreipiamas į /game/:id
```

### Error Flow 1 - Jau tame žaidime:
```
1. Vartotojas prisijungęs prie žaidimo A
2. Bandai prisijungti prie žaidimo A dar kartą
3. ❌ Error: "already in this game"
4. Frontend turėtų nukreipti į /game/A vietoj rodyti error
```

### Error Flow 2 - Jau kitame žaidime:
```
1. Vartotojas prisijungęs prie žaidimo A
2. Bando prisijungti prie žaidimo B
3. ❌ Error: "already in another active game"
4. Frontend rodo pranešimą: "You must leave your current game first"
```

---

## 🎯 Frontend Improvements (TODO)

### 1. Disable Join Button jei vartotojas jau žaidime
```javascript
// Dashboard.jsx
const [userActiveGame, setUserActiveGame] = useState(null);

useEffect(() => {
  // Check if user is in any active game
  axios.get('/api/v1/games/my-active-game')
    .then(res => setUserActiveGame(res.data))
    .catch(() => setUserActiveGame(null));
}, []);

// Render:
<button 
  onClick={() => joinGame(game.room_code)}
  disabled={userActiveGame && userActiveGame.id !== game.id}
>
  {userActiveGame && userActiveGame.id !== game.id 
    ? "Already in another game"
    : "Join Game"}
</button>
```

### 2. Rodyti "Resume Game" vietoj "Join" jei vartotojas jau tame žaidime
```javascript
{userActiveGame && userActiveGame.id === game.id ? (
  <button onClick={() => navigate(`/game/${game.id}`)}>
    Resume Game
  </button>
) : (
  <button onClick={() => joinGame(game.room_code)}>
    Join Game
  </button>
)}
```

### 3. Pridėti "Leave Game" funkcionalumą
```javascript
const leaveGame = async (gameId) => {
  try {
    await axios.post(`/api/v1/games/${gameId}/leave`);
    setUserActiveGame(null);
    fetchGames();
  } catch (err) {
    console.error('Failed to leave game:', err);
  }
};
```

---

## 📝 Backend Endpoint Suggestions (Future)

### 1. GET /api/v1/games/my-active-game
```go
// GetMyActiveGame returns the active game for current user
func (h *GameHandler) GetMyActiveGame(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	
	activeGame, err := h.gameService.GetUserActiveGame(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	if activeGame == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not in any active game"})
		return
	}
	
	c.JSON(http.StatusOK, activeGame)
}
```

### 2. POST /api/v1/games/:id/leave
```go
// LeaveGame removes player from a game
func (h *GameHandler) LeaveGame(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	gameID, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.gameService.LeaveGame(gameID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "left game successfully"})
}
```

---

## ✅ Testing Checklist

Po backend rebuild ir restart:

- [ ] Sukurti naują žaidimą A
- [ ] Bandyti join prie žaidimo A dar kartą → Should get "already in this game"
- [ ] Sukurti naują žaidimą B (kitas vartotojas)
- [ ] Bandyti join prie žaidimo B → Should get "already in another active game"
- [ ] Backend logs turėtų rodyti aiškius error message

---

## 🚀 Deployment

### 1. Rebuild backend:
```powershell
cd "c:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
docker compose build app --no-cache
```

### 2. Restart konteinerius:
```powershell
docker compose down
docker compose up -d
```

### 3. Test:
- Sukurti žaidimą
- Bandyti join prie kito žaidimo
- Tikėtis aiškesnio error message

---

**Data:** 2026-01-18  
**Status:** ✅ PATAISYTA  
**Files Changed:** 
- `internal/repository/game_player_repository.go` (pridėtas GetUserActiveGame)
- `internal/service/game_service.go` (atnaujintas JoinGame)


