package cards

import (
	"fmt"
	"math/rand"
)

// Suit represents a card suit
type Suit string

const (
	Hearts   Suit = "hearts"
	Diamonds Suit = "diamonds"
	Clubs    Suit = "clubs"
	Spades   Suit = "spades"
)

// Rank represents a card rank
type Rank string

const (
	Ace   Rank = "A"
	Two   Rank = "2"
	Three Rank = "3"
	Four  Rank = "4"
	Five  Rank = "5"
	Six   Rank = "6"
	Seven Rank = "7"
	Eight Rank = "8"
	Nine  Rank = "9"
	Ten   Rank = "10"
	Jack  Rank = "J"
	Queen Rank = "Q"
	King  Rank = "K"
)

// Card represents a playing card
type Card struct {
	Suit  Suit `json:"suit"`
	Rank  Rank `json:"rank"`
	Value int  `json:"value"` // Numeric value for ranking (2-14, Ace=14)
}

// Deck represents a deck of cards
type Deck struct {
	Cards []Card
}

// RankToValue converts a rank string to its numeric value
func RankToValue(rank Rank) int {
	switch rank {
	case Ace:
		return 14
	case King:
		return 13
	case Queen:
		return 12
	case Jack:
		return 11
	case Ten:
		return 10
	case Nine:
		return 9
	case Eight:
		return 8
	case Seven:
		return 7
	case Six:
		return 6
	case Five:
		return 5
	case Four:
		return 4
	case Three:
		return 3
	case Two:
		return 2
	default:
		return 0
	}
}

// NewCard creates a new card
func NewCard(suit Suit, rank Rank) Card {
	return Card{
		Suit:  suit,
		Rank:  rank,
		Value: RankToValue(rank),
	}
}

// NewDeck creates a full 52-card deck
func NewDeck() *Deck {
	suits := []Suit{Hearts, Diamonds, Clubs, Spades}
	ranks := []Rank{Ace, King, Queen, Jack, Ten, Nine, Eight, Seven, Six, Five, Four, Three, Two}

	deck := &Deck{
		Cards: make([]Card, 0, 52),
	}

	for _, suit := range suits {
		for _, rank := range ranks {
			deck.Cards = append(deck.Cards, NewCard(suit, rank))
		}
	}

	return deck
}

// Shuffle shuffles the deck using Fisher-Yates algorithm
func (d *Deck) Shuffle() {
	for i := len(d.Cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
}

// Draw removes and returns the top card from the deck
func (d *Deck) Draw() *Card {
	if len(d.Cards) == 0 {
		return nil
	}
	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return &card
}

// IsOnePlus checks if this card's rank is one higher than another card (with cyclic wrap: A+1=2)
func (c Card) IsOnePlus(other Card) bool {
	// Regular progression: 2->3, 3->4, ..., K->A
	if c.Value == other.Value+1 {
		return true
	}
	// Cyclic wrap: A (value 14) + 1 = 2 (value 2)
	if other.Value == 14 && c.Value == 2 {
		return true
	}
	return false
}

// IsSameSuit checks if this card has the same suit as another card
func (c Card) IsSameSuit(other Card) bool {
	return c.Suit == other.Suit
}

// IsHigherRank checks if this card's rank is higher than another card
func (c Card) IsHigherRank(other Card) bool {
	return c.Value > other.Value
}

// IsTrump checks if this card is a trump card
func (c Card) IsTrump(trumpSuit *string) bool {
	if trumpSuit == nil {
		return false
	}
	return string(c.Suit) == *trumpSuit
}

// String returns a string representation of the card
func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}

// FindNineOfSpades finds the index of 9 of spades in a hand
func FindNineOfSpades(hand []Card) int {
	for i, card := range hand {
		if card.Suit == Spades && card.Rank == Nine {
			return i
		}
	}
	return -1
}

// CardExists checks if a card exists in a hand
func CardExists(hand []Card, suit Suit, rank Rank) bool {
	for _, card := range hand {
		if card.Suit == suit && card.Rank == rank {
			return true
		}
	}
	return false
}

// RemoveCard removes a specific card from a hand
func RemoveCard(hand []Card, suit Suit, rank Rank) []Card {
	for i, card := range hand {
		if card.Suit == suit && card.Rank == rank {
			return append(hand[:i], hand[i+1:]...)
		}
	}
	return hand
}

// CountByRank counts cards of a specific rank in a hand
func CountByRank(hand []Card, rank Rank) int {
	count := 0
	for _, card := range hand {
		if card.Rank == rank {
			count++
		}
	}
	return count
}

// CountBySuit counts cards of a specific suit in a hand
func CountBySuit(hand []Card, suit Suit) int {
	count := 0
	for _, card := range hand {
		if card.Suit == suit {
			count++
		}
	}
	return count
}
