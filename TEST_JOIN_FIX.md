# 🧪 Join Game Fix - Testavimo Instrukcijos

## ✅ Kas Buvo Pataisyta

**Problema:** "Failed to join game: already in game" error net kai vartotojas nebandė join to paties žaidimo.

**Sprendimas:** Dabar sistema tikrina ar vartotojas yra **kitame** aktyviam žaidime ir grąžina aiškesnius error messages.

---

## 🔧 Rebuild ir Restart

### 1. Build'inas veikia fone

Backend Docker image dabar perkompiliuojamas su pataisymu.

### 2. Laukite pabaigos

Terminal parodys kai baigsis. Tai gali užtrukti 1-2 minutes.

### 3. Restart konteinerius

```powershell
cd "c:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
docker compose down
docker compose up -d
```

---

## 🧪 Kaip Testuoti

### Test 1: Join prie to paties žaidimo (savo)

**Žingsniai:**
1. Atidarykite http://localhost:3000/dashboard
2. Paspauskite "Create New Game"
3. Sukuriamas žaidimas A (jūs automatiškai pridedami kaip žaidėjas)
4. **Bandykite** join prie žaidimo A dar kartą

**Laukiamas rezultatas:**
```
❌ Error: "already in this game"
```

**Pastaba:** Frontend **turėtų** rodyti "Resume" mygtuką vietoj "Join", bet tai dar neįdiegta.

---

### Test 2: Join prie kito žaidimo

**Žingsniai:**
1. Esate jau žaidime A
2. Kitas vartotojas sukuria žaidimą B
3. Jūs bandote join prie žaidimo B

**Laukiamas rezultatas:**
```
❌ Error: "already in another active game"
```

**Tai yra teisingas behavior** - vartotojas negali būti dviejuose žaidimuose vienu metu.

---

### Test 3: Join po žaidimo pabaigos

**Žingsniai:**
1. Esate žaidime A
2. Žaidimas A baigiasi (state = 'finished')
3. Bandote join prie žaidimo B

**Laukiamas rezultatas:**
```
✅ Success - prisijungėte prie žaidimo B
```

**Pastaba:** Tai veiks, nes GetUserActiveGame tikrina tik aktyvius žaidimus (state IN ('waiting', 'playing')).

---

### Test 4: Normalus join (pirmas kartas)

**Žingsniai:**
1. Nesate jokiame žaidime
2. Kitas vartotojas sukuria žaidimą
3. Jūs join prie to žaidimo

**Laukiamas rezultatas:**
```
✅ Success - prisijungėte prie žaidimo
✅ Nukreipta į /game/:id
```

---

## 📊 Error Messages Comparison

### Prieš pataisymą:
```
❌ "already in game"
```
- Neaišku ar tai tas pats žaidimas ar kitas
- Neaišku kaip išspręsti problemą

### Po pataisymo:
```
❌ "already in this game" 
   (jei bandai join to paties žaidimo)

❌ "already in another active game"
   (jei bandai join kito žaidimo būdamas aktyviam žaidime)
```
- Aiškūs error messages
- Vartotojas žino kas negerai
- Frontend gali reaguoti atitinkamai

---

## 🔍 Debug su Backend Logs

Norint pamatyti ką backend'as daro:

```powershell
docker compose logs app -f
```

**Stebėkite:**
```
# Kai sukuriate žaidimą:
2026/01/18 ... Created game with ID: 1, Room Code: ABC123
2026/01/18 ... Added player (user_id=1) to game 1 at position 0

# Kai bandote join prie to paties:
2026/01/18 ... User 1 attempted to join game 1
2026/01/18 ... User already in this game (game_id=1, user_id=1)

# Kai bandote join prie kito:
2026/01/18 ... User 1 attempted to join game 2
2026/01/18 ... User has active game: game_id=1
2026/01/18 ... Cannot join game 2 - already in another active game
```

---

## 💡 Frontend Improvements (TODO)

Šie pakeitimai padėtų išvengti error'ų:

### 1. Disable "Join" mygtuką

```javascript
// Jei vartotojas jau žaidime:
<button disabled={userInGame}>
  {userInGame ? "Already in a game" : "Join Game"}
</button>
```

### 2. Rodyti "Resume" vietoj "Join"

```javascript
// Jei vartotojas tame žaidime:
{isMyGame ? (
  <button onClick={() => navigate(`/game/${game.id}`)}>
    Resume Game
  </button>
) : (
  <button onClick={() => joinGame(game.room_code)}>
    Join Game
  </button>
)}
```

### 3. "Leave Game" funkcionalumas

```javascript
// Jei vartotojas nori prisijungti prie kito:
<button onClick={leaveCurrentGame}>
  Leave Current Game
</button>
```

---

## 📝 Checklist Po Rebuild

- [ ] Rebuild baigėsi be klaidų
- [ ] Konteineriai restart'inti
- [ ] Visi 3 konteineriai veikia (postgres, app, frontend)
- [ ] Galite prisijungti per Google OAuth
- [ ] Dashboard rodomas
- [ ] Test 1: Join to paties žaidimo → Error "already in this game"
- [ ] Test 2: Join kito žaidimo → Error "already in another active game"
- [ ] Test 4: Normalus join veikia ✅

---

## 🎯 Laukiami Rezultatai

### Sėkminga situacija:
- ✅ Naujas vartotojas gali join prie laukiančio žaidimo
- ✅ Aiškūs error messages kai negali join
- ✅ Sistema neleidžia būti dviejuose žaidimuose

### Tikėtini errors (teisingi):
- "already in this game" - jei join to paties
- "already in another active game" - jei join kito būdamas žaidime
- "game already started" - jei žaidimas jau pradėtas
- "game is full" - jei žaidimas pilnas

---

## 🚀 Sekantys Žingsniai

1. **Laukti rebuild pabaigos**
2. **Restart konteinerius:** `docker compose down && docker compose up -d`
3. **Testuoti** visus 4 scenarijus
4. **Patikrinti** backend logs
5. **Jei veikia** - puiku! ✅
6. **Jei ne** - pažiūrėti logs ir debug

---

**Data:** 2026-01-18  
**Status:** 🔄 REBUILD IN PROGRESS  
**Laukiama:** Backend rebuild + restart + testing


