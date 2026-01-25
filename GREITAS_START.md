# 🎴 Fasiolas Card Game - Greitas Paleidimas

## ⚡ Greičiausias būdas paleisti

```powershell
# Eikite į projekto aplanką
cd "c:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"

# Paleiskite viską
.\start-all.ps1
```

Šis script'as:
- ✅ Sustabdo visus konteinerius
- ✅ Paleidžia iš naujo (PostgreSQL, Backend, Frontend)
- ✅ Pritaiko DB migracijas
- ✅ Perkrauna backend
- ✅ Patikrina, kad viskas veikia

---

## 🔧 Jei matote seną frontend versiją

```powershell
.\rebuild-frontend.ps1
```

Šis script'as:
- ✅ Sustabdo frontend
- ✅ Perkompiliuoja be cache
- ✅ Paleidžia naują versiją

---

## 🧪 Kaip testuoti OAuth

1. **Uždarykite** visus naršyklės langus
2. **Atidarykite** naują **Incognito** langą (`Ctrl+Shift+N`)
3. **Eikite** į: `http://localhost:3000/login`
4. **Paspauskite** "Google Login"
5. **Prisijunkite** su Google paskyra
6. ✅ **Turėtumėte** būti nukreipti į dashboard

---

## 📊 Ką daryti jei matote problemas

### Problema: "relation users does not exist"
**Sprendimas:** Migracijos nepritaikytos. Paleiskite:
```powershell
.\start-all.ps1
```

### Problema: Matote seną puslapio versiją
**Sprendimas:** Frontend cache. Paleiskite:
```powershell
.\rebuild-frontend.ps1
```

### Problema: Matote jau prisijungtą profilį incognito
**Sprendimas:** Frontend build'as senas. Paleiskite:
```powershell
.\rebuild-frontend.ps1
```
Po to:
- Uždarykite visus naršyklės langus
- Atidarykite naują incognito langą
- Eikite į http://localhost:3000/login

### Problema: Konteineriai nestartuoja
**Sprendimas:** Paleiskite viską iš naujo:
```powershell
docker compose down -v
docker compose up -d
```

---

## 📖 Išsamesni nurodymai

Žiūrėkite [KAIP_PALEISTI.md](KAIP_PALEISTI.md) - ten rasite detalius nurodymus.

---

## ✅ Kas buvo sutaisyta

1. ✅ **OAuth state parameter** - dabar tinkamai validuojamas
2. ✅ **DB migracijos** - visos lentelės sukurtos
3. ✅ **Docker konfigūracija** - sutvarkyti env vars ir paths
4. ✅ **Frontend/Backend susikalbėjimas** - veikia teisingai

---

## 🎯 Sistema veikia jei

- ✅ Visi 3 konteineriai paleisti (postgres, app, frontend)
- ✅ Backend grąžina OAuth URL su state parametru
- ✅ Login puslapis rodo Google/GitHub/Discord mygtukus
- ✅ Po prisijungimo nukreipia į dashboard

---

## 🚀 Sukurta

**Data:** 2026-01-18
**Status:** ✅ Veikia
**OAuth:** ✅ Google OAuth veikia
**Migracijos:** ✅ Pritaikytos


