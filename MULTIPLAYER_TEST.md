# 🎮 MULTIPLAYER TESTING - Quick Guide

## ✅ Sistema dabar leidžia keliems žaidėjams!

Konteineriai jau restart'inti su pataisymu.

---

## 🧪 Kaip Testuoti (2 žaidėjai viename žaidime)

### Žaidėjas 1 - Main Browser

1. **Atidarykite:** http://localhost:3000/login
2. **Prisijunkite** per Google (jei jau niste - OK)
3. **Eikite į:** Dashboard
4. **Paspauskite:** "Create New Game"
5. **Sukuriamas žaidimas** - jūs automatiškai pridedami

**Žaidimas sukurtas! Pamatykite room code - pvz: "ABC123"**

---

### Žaidėjas 2 - Incognito Browser

1. **Atidarykite** naują **Incognito** langą (Ctrl+Shift+N)
2. **Eikite į:** http://localhost:3000/login
3. **Prisijunkite** su **KITA** Google paskyra
   - Arba naudokite skirtingą Google account
4. **Eikite į:** Dashboard
5. **Paieškokite** žaidimo su room code iš Žaidėjo 1
6. **Paspauskite:** "Join Game"

---

## ✅ Laukiamas Rezultatas

### Žaidėjas 1:
- ✅ Sukūrė žaidimą
- ✅ Matote save kaip "Player 1"

### Žaidėjas 2:
- ✅ Prisijungė prie žaidimo
- ✅ Matote save kaip "Player 2"
- ✅ Matote Žaidėjo 1 informaciją

### Bendrai:
- ✅ Žaidime dabar 2 žaidėjai
- ✅ Abu matote vienas kitą
- ✅ Galite pradėti žaidimą!

---

## ❌ Jei Vis Dar Matote Error

**Error:** "already in this game"

**Sprendimas:** Tai reiškia, kad bandote join du kartus su **ta pačia** paskyra

**Ką daryti:**
1. Naudokite **skirtingą** Google account (Žaidėjas 2)
2. Arba atidarykite naują Incognito langą
3. Arba naudokite mobiliųjį telefoną

---

## 📊 Checklist

- [ ] Konteineriai veikia (docker ps)
- [ ] Žaidėjas 1 sukūrė žaidimą
- [ ] Žaidėjas 2 prisijungė prie žaidimo
- [ ] Abu matote vienas kitą
- [ ] Žaidime 2 žaidėjai
- [ ] ✅ Veikia!

---

## 🚀 Sekantis Žingsnis

Jei multiplayer veikia, galite:
- Pradėti žaidimą (tik creator gali)
- Dėti kortas
- Skambinti cheat
- Žaisti!

---

**STATUS: READY FOR MULTIPLAYER!** 🎮✅


