# 🎮 KAIP PAMATYTI START GAME MYGTUKĄ - INSTRUKCIJA

## ✅ Ką padariau

1. ✓ Pataisiau `Game.jsx` - **mygtukas visada rodomas** (žalia, didelis, su emoji)
2. ✓ Pataisiau `Dashboard.jsx` - game creator dabar **nukreipiamas į game room**
3. ✓ Pataisiau `package.json` - pridėjau `"dev"` script'ą
4. ✓ Paleisiau **backend** ir **frontend** serverius

---

## 🚀 TESTUOTI DABAR

### 1. Atidaryk INCOGNITO naršyklės langą
- Chrome/Edge: **Ctrl + Shift + N**

### 2. Eik į
```
http://localhost:3000
```

### 3. Prisijunk su Google

### 4. **SKIRTINGAI TESTUOK ABI SITUACIJAS:**

#### A) GAME CREATOR (žaidimo kūrėjas)
1. Spausk **"Create Game"** mygtuką
2. Pasirink žaidėjų skaičių (pvz. 2)
3. Spausk **"Create"**
4. **TURI BŪTI AUTOMATIŠKAI NUKREIPTAS Į GAME ROOM!**
5. Turėtum matyti:
   - Žalią mygtuką: **"🎮 Start Game"** (disabled, nes tik 1/2 žaidėjai)
   - Tekstą: **"✓ Logged in as: [vardas]"**

#### B) KITAS ŽAIDĖJAS (tavo draugas)
1. Nuo pirmojo žaidėjo paklaus **room code** (pvz. E30LR8)
2. Įvesti kodą į "Join by Room Code" input'ą
3. Spausk **"Join Game"**
4. **TURI BŪTI AUTOMATIŠKAI NUKREIPTAS Į GAME ROOM!**

#### C) GAME ROOM (jei abu žaidėjai yra)
1. **Abu matote:**
   - "Players Joined: 2/2"
   - Abu žaidėjai sąraše
   - **ŽALIĄ MYGTUKĄ: "🎮 Start Game"** (AKTYVUS - galite spausti!)

2. **Spauskit "Start Game"** - žaidimas turėtų prasidėti

---

## 🐛 JEI NEMATAI MYGTUKO

### Debug žingsniai:
1. Spausk **F12** (Developer Tools)
2. Eik į **Console** tab
3. Turėtum matyti:
   ```
   Current user: {id: 1, username: "...", ...}
   Players: [{...}, {...}]
   Is player in game: true/false
   ```

4. **Padaryk SCREENSHOT** ir parodyk man

### Jei Frontend nenaujinatasi:
1. Spausk **Ctrl + Shift + Delete** (Clear browsing data)
2. Pasirink "Cached images and files"
3. Spausk "Clear data"
4. **Atidaryk NAUJĄ Incognito langą**
5. Eik į `http://localhost:3000`

---

## 🎯 TIKSLAI

✅ **Game creator** → Create game → **auto-redirect to lobby**
✅ **Other player** → Join by code → **auto-redirect to lobby**
✅ **Both in lobby** → matoti **START GAME** mygtuką
✅ **Console** → matoti debug info

---

## ❓ PROBLEMOS APTARTOS

### Problema 1: "already in this game"
- **Sprendimas**: Game creator dabar automatiškai nukreipiamas į game room po create

### Problema 2: Frontend neatnaujina
- **Sprendimas**: Naudojame `npm start` (Create React App), ne Vite
- **Dev script** pridėtas į `package.json`

### Problema 3: Start button nematomas
- **Sprendimas**: Button nuo sąlyginiu render'avimo, visada rodomas
- **Pamatai lygiausiai po failų perkompiliacijo**

---

## 📋 SERVERIŲ STATUSAS

- ✓ Backend: `http://localhost:8080`
- ✓ Frontend: `http://localhost:3000` (startuoja ~10-15 sec)
- ✓ Database: lokali (development mode)

**Palaukk 10 sekundžių, kad frontend pasiruoš!**

---

Pasakyk, ką matai! 🎮

