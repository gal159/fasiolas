package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ExternalAPIService handles external API integration
type ExternalAPIService struct {
	apiURL string
	client *http.Client
}

// NewExternalAPIService creates a new external API service
func NewExternalAPIService(apiURL string) *ExternalAPIService {
	return &ExternalAPIService{
		apiURL: apiURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CardImage represents a card image from the API
type CardImage struct {
	Code   string `json:"code"`
	Image  string `json:"image"`
	Images struct {
		SVG string `json:"svg"`
		PNG string `json:"png"`
	} `json:"images"`
	Value string `json:"value"`
	Suit  string `json:"suit"`
}

// GetCardImage fetches card image from external API
func (s *ExternalAPIService) GetCardImage(suit, rank string) (*CardImage, error) {
	// Convert suit and rank to API format
	code := s.convertToCode(suit, rank)

	url := fmt.Sprintf("%s/deck/new/draw/?cards=%s", s.apiURL, code)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch card image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response struct {
		Success bool        `json:"success"`
		Cards   []CardImage `json:"cards"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success || len(response.Cards) == 0 {
		return nil, fmt.Errorf("no card found")
	}

	return &response.Cards[0], nil
}

// GetMultipleCardImages fetches multiple card images
func (s *ExternalAPIService) GetMultipleCardImages(cards []string) ([]CardImage, error) {
	if len(cards) == 0 {
		return []CardImage{}, nil
	}

	// Join card codes
	cardCodes := ""
	for i, code := range cards {
		if i > 0 {
			cardCodes += ","
		}
		cardCodes += code
	}

	url := fmt.Sprintf("%s/deck/new/draw/?cards=%s", s.apiURL, cardCodes)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch card images: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response struct {
		Success bool        `json:"success"`
		Cards   []CardImage `json:"cards"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("API request failed")
	}

	return response.Cards, nil
}

// convertToCode converts internal suit/rank to API code
func (s *ExternalAPIService) convertToCode(suit, rank string) string {
	// Rank conversion
	rankCode := rank
	switch rank {
	case "J":
		rankCode = "J"
	case "Q":
		rankCode = "Q"
	case "K":
		rankCode = "K"
	case "A":
		rankCode = "A"
	default:
		rankCode = rank
	}

	// Suit conversion
	suitCode := ""
	switch suit {
	case "hearts":
		suitCode = "H"
	case "diamonds":
		suitCode = "D"
	case "clubs":
		suitCode = "C"
	case "spades":
		suitCode = "S"
	}

	return rankCode + suitCode
}

// ShuffleDeck creates a new shuffled deck via API
func (s *ExternalAPIService) ShuffleDeck() (string, error) {
	url := fmt.Sprintf("%s/deck/new/shuffle/?deck_count=1", s.apiURL)

	resp, err := s.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to shuffle deck: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var response struct {
		Success bool   `json:"success"`
		DeckID  string `json:"deck_id"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success {
		return "", fmt.Errorf("failed to create deck")
	}

	return response.DeckID, nil
}

// CardStats represents statistics about card games
type CardStats struct {
	TotalGames      int     `json:"total_games"`
	AverageGameTime float64 `json:"average_game_time"`
	PopularCard     string  `json:"popular_card"`
}

// GetGameStats returns mock game statistics (example of external API integration)
func (s *ExternalAPIService) GetGameStats() (*CardStats, error) {
	// This is a mock implementation
	// In a real scenario, this would call an actual stats API
	return &CardStats{
		TotalGames:      1000,
		AverageGameTime: 15.5,
		PopularCard:     "AS", // Ace of Spades
	}, nil
}
