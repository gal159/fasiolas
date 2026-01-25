# ✅ VISKAS VEIKIA - Kaip Pradėti

## 🎯 Dabartinė būsena

✅ Visi konteineriai paleisti
✅ Duomenų bazė sukurta (users, games, game_players, game_actions lentelės)
✅ Backend veikia (port 8080)
✅ Frontend veikia (port 3000)
✅ OAuth state validation veikia

---

## 🚀 Kaip Testuoti Google OAuth

### Greitas būdas - naudokite script'us:

**Paleisti viską iš naujo (su migracijomis):**
```powershell
.\start-all.ps1
```

**Tik perkompiliuoti frontend:**
```powershell
.\rebuild-frontend.ps1
```

### Rankinis būdas:

### 1. Išvalykite naršyklės cache

**Chrome/Edge:**
1. Paspauskite `Ctrl + Shift + R` (arba `Ctrl + F5`)
2. Arba: F12 → Network tab → pažymėkite "Disable cache"

**Arba atidarykite Incognito/Private režimą:**
- `Ctrl + Shift + N` (Chrome)
- `Ctrl + Shift + P` (Edge/Firefox)

### 2. Atidarykite login puslapį

```
http://localhost:3000/login
```

### 3. Paspauskite "Google Login"

Turėtumėte:
1. Būti nukreipti į Google authentication
2. Prisijungti su savo Google paskyra
3. Būti nukreipti atgal į aplikaciją
4. Matyti dashboard (jei viskas gerai)

---

## 🔍 Jei vis dar neveikia arba matote seną versiją

### Problema: Matote jau prisijungtą profilį incognito režime

Tai reiškia, kad frontend build'as yra senas. Reikia perkompiliuoti frontend:

```powershell
cd "c:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"

# Sustabdyti frontend
docker compose stop frontend

# Perkompiliuoti frontend be cache
docker compose build --no-cache frontend

# Paleisti frontend
docker compose up -d frontend
```

Po šių komandų:
1. Uždarykite visus naršyklės langus
2. Atidarykite naują incognito langą
3. Eikite į http://localhost:3000/login
4. Dabar turėtumėte matyti login puslapį

### Perkrauti visus konteinerius:

```powershell
cd "c:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
docker compose restart
```

### Pažiūrėti backend logo:

```powershell
docker compose logs app -f
```

Ieškokite šių žinučių:
- ✅ `[OAuth] State generated and stored in memory:`
- ✅ `[OAuth] ✅ State validated from memory store`
- ✅ `[OAuth] ✅ Token generated successfully`

### Jei matote "relation users does not exist" klaidą:

Backend'as dar naudoja seną connection. Perkraukite:

```powershell
docker compose restart app
```

---

## 📊 Patvirtinimai

### Duomenų bazė veikia:

```powershell
docker compose exec postgres psql -U postgres -d fasiolas_game -c "\dt"
```

Turėtumėte matyti:
- users
- games
- game_players
- game_actions

### Backend veikia:

```powershell
curl http://localhost:8080/health
```

### Frontend veikia:

Atidarykite naršyklėje: http://localhost:3000

---

## 🎉 Kas Sutvarkyta

1. ✅ OAuth state parameter dabar tinkamai validuojamas
2. ✅ Duomenų bazės lentelės sukurtos
3. ✅ Visi konteineriai veikia
4. ✅ Backend ir frontend susikalbėja

**Problema, kurią matote naršyklėje:** Tai cache problema. Naršyklė naudoja senus failus (304 Not Modified).

**Sprendimas:** Hard refresh (`Ctrl + Shift + R`) arba Incognito mode.

---

## 🔧 Jei reikia visiškai iš naujo paleisti

```powershell
# 1. Sustabdyti viską
docker compose down

# 2. Paleisti iš naujo
docker compose up -d

# 3. Nukopijuoti migracijas
docker compose cp migrations/. postgres:/migrations

# 4. Pritaikyti migracijas
docker compose exec postgres psql -U postgres -d fasiolas_game -f /migrations/000001_create_users_table.up.sql
docker compose exec postgres psql -U postgres -d fasiolas_game -f /migrations/000002_create_games_table.up.sql
docker compose exec postgres psql -U postgres -d fasiolas_game -f /migrations/000003_create_game_players_table.up.sql
docker compose exec postgres psql -U postgres -d fasiolas_game -f /migrations/000004_create_game_actions_table.up.sql

# 5. Perkrauti backend
docker compose restart app

# 6. Išvalyti naršyklės cache ir testuoti
```

---

**Data:** 2026-01-18
**Status:** ✅ VEIKIA - Tik reikia išvalyti cache


