# Žaidimo Taisyklės - Implementacija

## Bendrosios Nuostatos

### Žaidėjų skaičius
- **2–6 žaidėjai** (rekomenduojama 4)

### Kaladė
- Standartinė **52 kortų kaladė** (be jokerių)

### Kortų Eiliškumas
```
2 < 3 < 4 < 5 < 6 < 7 < 8 < 9 < 10 < J < Q < K < A
```

### Mastys (rūšys)
```
♠ (Picos) - Black
♥ (Širdys) - Red
♦ (Deimantai) - Red
♣ (Krūmais) - Black
```

## Žaidimas vyksta dviem fazėmis

### FAZĖ 1: Kaupimas ir Baudų Mechanika
- Žaidėjai deda kortas vienas ant kito pagal **+1 ciklinę taisyklę**
- Pažaidimas vyksta iki kol baigsis kaladė

### FAZĖ 2: Strateginis Žaidimas
- Žaidėjai žaidžia strateginiu būdu su **Trump (Trimi) kostiumu**
- Žaidimas tęsiasi kol lieka tik vienas žaidėjas su kortomis

---

## FAZĖ 1: Kaupimas ir Baudų Mechanika

### Paruošimas žaidimui

1. **Kaladė sumaišoma**
2. **Kiekvienas žaidėjas gauna po 1 kortą** - ji padedama atversta (tai yra jo krūvos viršus)
3. **Pirmas žaidėjas** - tas, kurio viršutinė korta turi **mažiausią rangą**
4. **Likusios kortos** - sudaro traukiamąją kaladę

### +1 Ciklinė Taisyklė (Apibrėžimas)

Kortą galima dėti ant krūvos, jei jos **rangas yra +1** nuo tos krūvos viršutinės kortos rango, **cikliškai**:

```
5 → 6 (standartinis)
K → A (standartinis)
A → 2 (CIKLINIS PERSIPLOJIMAS!)
2 → 3 (standartinis)
```

**Svarbu**: A (Tūzas) yra didžiausia korta, bet A+1 = 2 (mažiausia).

---

## Žaidėjo Ėjimas (Phase 1)

### A) Viršūnės dėjimas (prieš traukiant)

**Žaidėjo ėjimas visada prasideda nuo bandymo padėti savo viršutinę kortą.**

#### A1) Jei savo viršutinę kortą gali padėti ant **bent vieno KITO žaidėjo krūvos** (+1 taisyklė):
- **PRIVALAI ją padėti** ant vieno iš tų kitų žaidėjų
- **Po padėjimo** vėl patikrini savo naują viršutinę kortą
- **Kartoji žingsnį**, kol nebegali padėti ant kitų

#### A2) Jei NEGALI padėti viršūnės ant nė vieno kito žaidėjo, **bet GALI padėti ant SAVĘS** (+1 taisyklė ant savęs):
- Padedi viršutinę kortą ant savo krūvos (tai padidina tavo krūvą)
- **Po tokio padėjimo TU PRIVALAI TRAUKTI kortą iš kaladės** (pereini į B žingsnį)

#### A3) Jei NEGALI padėti viršūnės nei ant kitų, nei ant savęs:
- Tu privalai traukti kortą iš kaladės (pereini į B žingsnį)

### B) Kortos Traukimas (kai reikia traukti)

Kai pagal A žingsnį reikia traukti:

**Žaidėjas traukia 1 kortą iš kaladės** (tai „ištraukta korta").

Po traukimo patikrinama, **kur ją galima dėti** pagal +1 taisyklę:

#### B1) Jei ištraukta korta **tinka ant bent vieno kito žaidėjo**:
- **PRIVALAI padėti ją** ant vieno iš tų kitų žaidėjų
- **Ėjimas BAIGIASI**

#### B2) Jei ištraukta korta **netinka nė ant vieno kito žaidėjo, bet tinka ant savęs**:
- **Padedi ją ant savęs**
- **Ėjimas BAIGIASI**

#### B3) Jei ištraukta korta **netinka niekur** pagal +1 taisyklę:
- **Padedi ją ant savęs** (pasiimi sau)
- **Ėjimas BAIGIASI** - eilė pereina kitam žaidėjui

---

## Ėjimo Logika - Pseudokodas

```pseudo
PLAYER_TURN():
  // A) Try to place on others repeatedly
  LOOP:
    IF player.topCard can match ANY other player:
      Place on that player
      Update player.topCard to next card
      CONTINUE LOOP  // Check again
    ELSE IF player.topCard can match SELF:
      Place on self
      MUST DRAW (go to B)
      BREAK
    ELSE:
      MUST DRAW (go to B)
      BREAK
  
  // B) Draw and handle
  DRAW_PHASE():
    drawnCard = deck.pop()
    
    IF drawnCard can match ANY other player:
      Place on that player
      TURN ENDS
    ELSE IF drawnCard can match SELF:
      Place on self
      TURN ENDS
    ELSE:
      Place on self (keep card)
      TURN ENDS
  
  NEXT PLAYER
```

---

## Konkreti Implementacija Go'je

### Duomenų Struktūra

```go
type GamePlayer struct {
    Position   int           // Player position (0, 1, 2, ...)
    Cards      []Card        // All cards in player's pile (ordered)
    TopCard    *Card         // The top card (visible) - same as Cards[len-1]
    CardCount  int           // Number of cards
}

type Game struct {
    Phase                 int     // 1 or 2
    CurrentPlayerPosition *int    // Which position plays now
    DeckCards            []Card   // Remaining deck
    TableCards           []Card   // Phase 2 only
}
```

### Funkcijos

#### 1. **CanPlaceCard** - Patikrina ar gali dėti
```go
func (e *Engine) CanPlaceCard(currentPlayer GamePlayer, allPlayers []GamePlayer) (bool, int, bool)
// Returns:
// - canPlace: true/false
// - targetPos: position of target player (-1 if can't place)
// - isOnSelf: true if can only place on self
```

#### 2. **ExecutePhase1Turn** - Pilnas ėjimas
```go
func (e *Engine) ExecutePhase1Turn(game *Game, players []GamePlayer) (*Phase1TurnState, error)
// Executes complete A + B phase
// Returns action taken: "place_on_other", "place_on_self_then_draw", etc.
```

#### 3. **DrawPhase1Card** - Traukimas ir automatinis dėjimas
```go
func (e *Engine) DrawPhase1Card(game *Game, player *GamePlayer, allPlayers []GamePlayer) (*Phase1TurnState, error)
// B1, B2, B3 logic
```

### Ciklinis IsOnePlus Čekis

```go
func (c Card) IsOnePlus(other Card) bool {
    // Standard: 5+1=6, K+1=A
    if c.Value == other.Value+1 {
        return true
    }
    // Cyclic: A (14) + 1 = 2
    if other.Value == 14 && c.Value == 2 {
        return true
    }
    return false
}
```

---

## Fazės Sąlyga

### Kada baigiasi Fazė 1?
- Kai **kaladė (DeckCards) yra tuščia** ir žaidėjas negalėjo dėti

### Kada prasideda Fazė 2?
- Kada baigiasi Fazė 1
- **Trumpas kostiumus** nustatomas iš paskutinio žaidėjo krūvos viršaus (nevadintas Pikų Devynetas)

---

## Svarbi Pastaba: Viršūnės Koncepcija

- **Viršūnė (TopCard)** = poslednia korta žaidėjo krūvoje = `Cards[len(Cards)-1]`
- Kai žaidėjas deda kortą, **nauja viršūnė** = ankstesnė korta
- **Apatinė korta** = kur bus dedamas + 1

### Pavyzdys
```
Player's pile (bottom to top):
[ 5, 8, 9 ]
            ↑ TopCard = 9

Can place 10? YES (9+1=10)
After placing 10 on another player:
[ 5, 8 ]
      ↑ New TopCard = 8
```

---

## Testas - Scenarijus

### Scenarijus 1: Tiesioginė Dėjimas Kitam
```
Player A top: 5♥
Player B top: 4♣
Player C top: Q♠

Player A's turn:
5♥ + 1 = 6
Gali dėti ant Player B? 5♥ is not 5? YES! (4+1=5, but we need 5♥ to be +1 to 4♣)
WAIT: 5 is not 4+1, so NO.
Actually: (4♣ value = 4, 5♥ value = 5, 5 = 4+1) YES!
→ Place 5♥ on Player B
→ Check new top (8♠): Can place on others? 8 on 5? NO. On 4? NO. On Q? NO.
→ Can place on self? No next card to compare... MUST DRAW
```

---

## Implementacijos Patikrinimas

Visas šis failas dokumentuoja taisykles, kurias turi sekti `internal/game/engine.go` failas.

**Pagrindinės Funkcijos:**
1. ✅ `IsOnePlus` - ciklinis rangas
2. ✅ `CanPlaceCard` - A1/A2/A3 logika
3. ✅ `ExecutePhase1Turn` - A fazė su loop'u
4. ✅ `DrawPhase1Card` - B fazė (B1/B2/B3)

**Testo Sutartis:**
- [ ] 2 žaidėjai, 52 kortos
- [ ] A→2 persiplojimas veikia
- [ ] Privalomas dėjimas ant kitų
- [ ] Savęs dėjimas + traukimas
- [ ] Fazė 1 baigiasi teisingai
- [ ] Fazė 2 prasideda teisingai


