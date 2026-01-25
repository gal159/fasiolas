# ✅ KAIP PAMATYTI "START GAME" MYGTUKĄ

## 🎯 Problema
Nerodomas "Start Game" mygtukas game lobby, nors yra 2 žaidėjai.

## 🔧 Ką padariau

1. **Pataisiau `Game.jsx` kodą** - mygtukas dabar VISADA rodomas (be jokių sąlygų)
2. **Mygtukas dabar turi:**
   - 🎮 Emoji ikoną
   - Didesnį fontą (text-xl)
   - Žalią spalvą
   - Hover efektą (scale-105)
   - Tekstą po mygtuku, rodantį ar esi prisijungęs

## 📋 Ką TU turi padaryti

### Žingsnis 1: Paleisti serverius
Jei serveriai dar neveikia, turėjo pasileisti automatiškai. 
Patikrink, ar matai naują PowerShell langą su "FRONTEND DEV SERVER".

### Žingsnis 2: Palaukti
Palaukk ~10-15 sekundžių, kol Vite dev serveris pasileidžia.
Naujame PowerShell lange turėtum matyti:
```
VITE v4.x.x  ready in xxx ms

➜  Local:   http://localhost:3000/
```

### Žingsnis 3: Atidaryti naršyklę TEISINGAI

**❌ NEDARYK:**
- Neatidaryk esamo lango
- Nenaudok paprastos naršyklės
- Nespausk F5 senoje kortelėje

**✅ DARYK:**
1. **UŽDARYK VISUS** naršyklės langus
2. Paspausk **Ctrl + Shift + N** (Chrome/Edge) - atsidarys INCOGNITO langas
3. Eik į: `http://localhost:3000`
4. Prisijunk su Google
5. Eik į Game Room (Game #9) arba sukurk naują

### Žingsnis 4: Patikrinti

Lobby puslapyje turėtum matyti:

```
Waiting for Players
┌─────────────────────────┐
│ Players Joined: 2/2     │
│ Room Code: E30LR8       │
└─────────────────────────┘

Players in Game:
- Džiugas Čiuta (Position 1)
- Džiugas Elite (Position 2)

[ 🎮 Start Game ]  ← DIDELIS ŽALIAS MYGTUKAS
✓ Logged in as: Džiugas Elite
```

## 🐛 Jei vis tiek nematai mygtuko

### Debug žingsniai:

1. **Paspausk F12** (Developer Tools)
2. **Eik į Console tab**
3. **Ieškok šių message'ų:**
   ```
   Current user: {id: 1, username: "...", ...}
   Players: [{...}, {...}]
   Is player in game: true/false
   ```

4. **Padaryk screenshot Console** ir parodyk man

5. **Patikrink Network tab:**
   - Ar yra `/api/v1/games/[id]` request'as?
   - Ar jis gražina duomenis?

## 🚀 Alternatyvus būdas

Jei visiškai neveikia dev mode, galima paleisti per Docker:

```powershell
cd "C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
docker compose down
docker compose up --build -d
```

Tada eik į: `http://localhost:3000`

## ❓ Vis dar neveikia?

Pasakyk man:
1. Ar matai console log'us? Jei taip, ką jie rodo?
2. Ar matai tekstą "✓ Logged in as: [vardas]"?
3. Ar matai bent kokį nors elementą po "Players in Game" sąrašu?
4. Ar Vite dev serveris tikrai veikia? (patikrink PowerShell langą)

---

## 📝 Techninis paaiškinimas

Kodas, kuris turėtų rodyti mygtuką (`Game.jsx`, line ~128):
```jsx
<div className="mt-8">
  <button
    onClick={handleStartGame}
    disabled={players.length < 2}
    className="bg-green-600 hover:bg-green-700 disabled:opacity-50 ..."
  >
    {players.length < 2 ? '⏳ Need at least 2 players' : '🎮 Start Game'}
  </button>
  <p className="text-gray-400 text-sm mt-2">
    {user ? `✓ Logged in as: ${user.username}` : '✗ Not logged in'}
  </p>
</div>
```

Mygtukas dabar NETURI jokių conditional render sąlygų (`{user && ...}` buvo pašalintas).
Jis turėtų būti rodomas VISADA.

Jei nematai - tai reiškia, kad:
1. Frontend dev serveris nerodo naujausios versijos (cache problema)
2. Arba component'as išvis nerenderinamas (tada nematytum ir "Players in Game")

