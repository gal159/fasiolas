# 🧪 Quick Test Guide - Card Placement Fix

## 🎮 Test Now!

**URL**: http://localhost:3000

### Test Case 1: Place Card on Opponent ✅

**Setup:**
1. Login and create a 2-player game
2. Start the game - both players get 1 card
3. Player 1 draws a card

**Action:**
1. Click on opponent's pile to place the card

**Expected Result:**
- ✅ Card is placed (no "invalid placement" error)
- ✅ Card moves from center to opponent's pile
- ✅ Check if +1 applies:
  - If YES → Your turn continues, draw button appears
  - If NO → Turn ends, next player's turn

### Test Case 2: Place Card on Yourself ✅

**Setup:**
1. Player draws a card

**Action:**
1. Click on your own pile to place the card

**Expected Result:**
- ✅ Card is placed (always succeeds)
- ✅ Card moves to your pile
- ✅ Check +1 rule:
  - If +1 → Turn continues
  - If no +1 → Turn ends

### Test Case 3: Turn Continuation Logic ✅

**Scenario A: +1 Applies**
```
Your card: 6♦
Opponent's top: 5♠
Place: 6 on 5? (6 = 5 + 1) YES ✅

Result: Turn CONTINUES
- Card placed ✅
- Deck button clickable ✅
- Still your turn ✅
```

**Scenario B: +1 Doesn't Apply**
```
Your card: 8♦
Opponent's top: 5♠
Place: 8 on 5? (8 ≠ 5 + 1) NO ❌

Result: Turn ENDS
- Card placed ✅
- Your turn changed ✅
- Next player's turn ✅
```

### Test Case 4: Multiple Placements (Chain) ✅

1. Draw card #1 that is +1 from opponent
2. Place on opponent → Turn continues
3. Draw card #2 that is +1 from another opponent
4. Place on that opponent → Turn continues
5. Draw card #3 that is NOT +1 from anyone
6. Place on yourself → Turn ends (or continues if +1 to your pile)

## ✅ Success Indicators

- [ ] Can place cards on opponent without error
- [ ] Can place cards on yourself without error
- [ ] Turn continues when +1 applies
- [ ] Turn ends when +1 doesn't apply
- [ ] No "invalid placement" error messages
- [ ] Game flows smoothly

## 🐛 If You See Errors

1. **"Invalid placement" still showing**
   - Hard refresh browser (Ctrl+F5)
   - Clear browser cache
   - Check server logs: `docker logs fasiolas_app`

2. **Turn not continuing when it should**
   - Check +1 rule is being evaluated correctly
   - Look at server logs for placement details

3. **Server not responding**
   - Check Docker: `docker ps`
   - Restart: `docker-compose restart`

## 📊 Card Values for +1 Testing

```
2  = 2
3  = 3
4  = 4
5  = 5
6  = 6
7  = 7
8  = 8
9  = 9
10 = 10
J  = 11
Q  = 12
K  = 13
A  = 14

+1 Rule Examples:
- 2 on A (2 on 14) → NO (2 ≠ 14 + 1)
- 3 on 2 (3 on 2) → YES (3 = 2 + 1)
- 4 on 3 (4 on 3) → YES (4 = 3 + 1)
- J on 10 (11 on 10) → YES (11 = 10 + 1)
- Q on J (12 on 11) → YES (12 = 11 + 1)
- K on Q (13 on 12) → YES (13 = 12 + 1)
- A on K (14 on 13) → YES (14 = 13 + 1)
- 2 on A (2 on 14) → YES (2 = 14 + 1, cyclic rule)
```

## 🎯 Expected Behavior Summary

| Card Placement | Result | Turn Continues? |
|---|---|---|
| Any card on opponent | ✅ Success | Check +1 |
| Any card on self | ✅ Success | Check +1 |
| +1 applies | ✅ Success | ✅ YES |
| +1 doesn't apply | ✅ Success | ❌ NO |
| Invalid target | ❌ Error | N/A |

---

**Status**: 🟢 Ready for Testing  
**All Systems**: ✅ Online
