# ✅ SISTEMA PARUOŠTA - Greitas Paleidimas

## 🚀 PowerShell Script'ai

Sukūriau 3 PowerShell script'us lengvam paleidimui:

### 1. `start-all.ps1` - Pilnas paleidimas

Paleidžia viską iš naujo su migracijomis:

```powershell
.\start-all.ps1
```

**Kas atliekama:**
- ✅ Sustabdo visus konteinerius
- ✅ Paleidžia PostgreSQL, Backend, Frontend
- ✅ Kopijuoja ir pritaiko DB migracijas
- ✅ Perkrauna backend
- ✅ Patikrina DB lenteles
- ✅ Rodo konteinerių būseną

**Kada naudoti:** Pirmas kartas arba po `docker compose down -v`

---

### 2. `quick-start.ps1` - Greitas paleidimas

Tiesiog paleidžia jau sukurtus konteinerius:

```powershell
.\quick-start.ps1
```

**Kas atliekama:**
- ✅ Paleidžia konteinerius
- ✅ Laukia 5 sek
- ✅ Rodo būseną

**Kada naudoti:** Kasdien, jei konteineriai jau sukurti

---

### 3. `rebuild-frontend.ps1` - Frontend perkompiliavimas

Perkompiliuoja frontend be cache:

```powershell
.\rebuild-frontend.ps1
```

**Kas atliekama:**
- ✅ Sustabdo frontend
- ✅ Perkompiliuoja be cache
- ✅ Paleidžia naują versiją

**Kada naudoti:** Kai matote seną frontend versiją arba cache problemas

---

## 📊 Kaip Naudoti

### Pirmas paleidimas:

```powershell
cd "c:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
.\start-all.ps1
```

### Kasdieninis paleidimas:

```powershell
cd "c:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
.\quick-start.ps1
```

### Jei matote seną versiją:

```powershell
.\rebuild-frontend.ps1
```

---

## 🧪 Testavimas

**Po bet kurio script'o paleidimo:**

1. **Uždarykite** visus naršyklės langus
2. **Atidarykite** naują **Incognito** langą (`Ctrl + Shift + N`)
3. **Eikite** į: `http://localhost:3000/login`
4. **Paspauskite** "Google Login"
5. **Prisijunkite** su Google paskyra
6. ✅ **Turėtumėte** būti nukreipti į dashboard

---

## ⚠️ Problemos ir Sprendimai

### Problema: "relation users does not exist"

**Sprendimas:**
```powershell
.\start-all.ps1
```
(Tai pritaikys migracijas)

---

### Problema: Matote seną puslapį

**Sprendimas:**
```powershell
.\rebuild-frontend.ps1
```
Po to:
- Uždarykite visus naršyklės langus
- Atidarykite Incognito režimu
- Eikite į http://localhost:3000/login

---

### Problema: Konteineriai nestartuoja

**Sprendimas:**
```powershell
docker compose down -v
.\start-all.ps1
```

---

### Problema: PowerShell klaidos

**Sprendimas:** Įsitikinkite kad esate projekto aplanke:
```powershell
cd "c:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
```

---

## ✅ Kas Veikia

- ✅ **Backend:** http://localhost:8080
- ✅ **Frontend:** http://localhost:3000
- ✅ **PostgreSQL:** localhost:5432
- ✅ **OAuth State Validation:** Veikia
- ✅ **DB Migracijos:** Pritaikytos

---

## 📖 Daugiau Informacijos

- **Detalūs nurodymai:** [KAIP_PALEISTI.md](KAIP_PALEISTI.md)
- **OAuth fix dokumentacija:** [START_HERE.md](START_HERE.md)
- **Vizualinės diagramos:** [FIX_VISUAL_DIAGRAMS.md](FIX_VISUAL_DIAGRAMS.md)

---

**Data:** 2026-01-18
**Status:** ✅ VEIKIA
**Script'ai sukurti:** ✅ start-all.ps1, quick-start.ps1, rebuild-frontend.ps1


