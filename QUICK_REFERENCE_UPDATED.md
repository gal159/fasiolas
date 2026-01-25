# 🎮 Quick Reference: Card Placement Rules (Updated)

## 🎯 The Golden Rules

### Rule 1: Placing on Opponents
```
✅ REQUIRES +1 rule validation
🔄 Turn CONTINUES (must draw again)
```

### Rule 2: Placing on Yourself
```
✅ ALWAYS ALLOWED (any card)
🔄 Turn continues IF +1 applies
❌ Turn ends IF +1 doesn't apply
```

---

## 📋 Quick Examples

### ✅ Valid: Place on Opponent
```
Your card: J♦ (11)
Opponent's top: 10♥ (10)
11 = 10 + 1? YES ✅
→ Place on opponent
→ Must draw again
```

### ❌ Invalid: Place on Opponent
```
Your card: J♦ (11)
Opponent's top: 8♦ (8)
11 = 8 + 1? NO ❌
→ Cannot place (Error)
```

### ✅ Valid: Place on Self (Turn Continues)
```
Your card: 9♥ (9)
Your top: 8♦ (8)
9 = 8 + 1? YES ✅
→ Place on self
→ Can draw again
```

### ✅ Valid: Place on Self (Turn Ends)
```
Your card: J♦ (11)
Your top: 8♦ (8)
11 = 8 + 1? NO
→ Place on self ✅
→ Turn ENDS
→ Next player
```

---

## 🎯 Decision Tree

```
Have a card to place?
│
├─ Want to place on OPPONENT?
│  │
│  ├─ Does +1 rule apply?
│  │  ├─ YES → ✅ Place it → Turn continues → Draw again
│  │  └─ NO  → ❌ ERROR "invalid placement"
│  │
│
└─ Want to place on YOURSELF?
   │
   └─ ALWAYS ALLOWED ✅
      │
      ├─ Does +1 rule apply?
      │  ├─ YES → Turn continues → Can draw again
      │  └─ NO  → Turn ENDS → Next player
```

---

## 📊 At a Glance

| Target | +1 Required? | Always Valid? | Turn Continues? |
|--------|--------------|---------------|-----------------|
| Opponent | ✅ YES | NO | ✅ YES (must draw) |
| Self (with +1) | NO | ✅ YES | ✅ YES (can draw) |
| Self (no +1) | NO | ✅ YES | ❌ NO (turn ends) |

---

## 🔍 Quick Checks

### ✅ Is Placement Valid?
```javascript
// Opponent?
if (target === opponent) {
  return card.value === opponent.topCard.value + 1;
}

// Self?
if (target === self) {
  return true; // Always valid!
}
```

### ✅ Does Turn Continue?
```javascript
// Check if +1 rule applies
const plusOneApplies = 
  card.value === target.topCard.value + 1;

if (target === self && !plusOneApplies) {
  return false; // Turn ends
}

return true; // Turn continues
```

---

## 🎮 Testing Checklist

- [ ] Place valid card on opponent → Turn continues ✅
- [ ] Place invalid card on opponent → Error ❌
- [ ] Place any card on self (with +1) → Turn continues ✅
- [ ] Place any card on self (no +1) → Turn ends ✅

---

## 🚀 Server Status

```bash
# Check containers
docker ps

# Expected output:
✅ fasiolas_frontend - Port 3000
✅ fasiolas_app      - Port 8080
✅ fasiolas_postgres - Port 5432
```

---

## 📞 Quick Links

- **Frontend**: http://localhost:3000
- **Backend**: http://localhost:8080
- **API Docs**: http://localhost:8080/swagger/index.html

---

## 🐛 Quick Debug

```bash
# Backend logs
docker logs fasiolas_app --tail 50

# Check game state
curl http://localhost:8080/api/v1/games/{id} \
  -H "Authorization: Bearer {token}"

# Look for errors
docker logs fasiolas_app 2>&1 | grep -i error
```

---

**Updated**: January 25, 2026  
**Status**: 🟢 Live and Ready
