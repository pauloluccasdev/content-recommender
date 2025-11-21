package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// RecommendationService define a interface para operações de recomendações
type RecommendationService interface {
	GetRecommendations(userID uint, topN int, method string) ([]uint, error)
	NotifyNewInteraction(userID, contentID uint, interactionType string, rating *float64) error
}

type recommendationService struct {
	recommenderURL string
	httpClient     *http.Client
}

// NewRecommendationService cria uma nova instância do RecommendationService
func NewRecommendationService() RecommendationService {
	recommenderURL := os.Getenv("RECOMMENDER_URL")
	if recommenderURL == "" {
		recommenderURL = "http://recommender:8000"
	}

	return &recommendationService{
		recommenderURL: recommenderURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// RecommendationRequest é o payload para o motor Python
type RecommendationRequest struct {
	UserID int    `json:"user_id"`
	TopN   int    `json:"top_n"`
	Method string `json:"method"`
}

// RecommendationResponse é a resposta do motor Python
type RecommendationResponse struct {
	UserID         int                      `json:"user_id"`
	Recommendations []ContentRecommendation `json:"recommendations"`
	Method         string                  `json:"method"`
}

// ContentRecommendation representa uma recomendação individual
type ContentRecommendation struct {
	ContentID int     `json:"content_id"`
	Score     float64 `json:"score"`
	Title     string  `json:"title"`
}

// GetRecommendations busca recomendações do motor Python
func (s *recommendationService) GetRecommendations(userID uint, topN int, method string) ([]uint, error) {
	if topN <= 0 || topN > 50 {
		topN = 10
	}
	if method == "" {
		method = "similarity"
	}

	reqBody := RecommendationRequest{
		UserID: int(userID),
		TopN:   topN,
		Method: method,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	url := fmt.Sprintf("%s/recommendations/", s.recommenderURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar motor de recomendação: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro do motor de recomendação (status %d): %s", resp.StatusCode, string(body))
	}

	var recommendationResp RecommendationResponse
	if err := json.NewDecoder(resp.Body).Decode(&recommendationResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	// Extrair apenas os IDs dos conteúdos recomendados
	contentIDs := make([]uint, len(recommendationResp.Recommendations))
	for i, rec := range recommendationResp.Recommendations {
		contentIDs[i] = uint(rec.ContentID)
	}

	return contentIDs, nil
}

// InteractionRequest é o payload para notificar nova interação
type InteractionRequest struct {
	UserID          int      `json:"user_id"`
	ContentID       int      `json:"content_id"`
	InteractionType string   `json:"interaction_type"`
	Rating          *float64 `json:"rating,omitempty"`
}

// NotifyNewInteraction notifica o motor Python sobre uma nova interação
func (s *recommendationService) NotifyNewInteraction(userID, contentID uint, interactionType string, rating *float64) error {
	reqBody := InteractionRequest{
		UserID:          int(userID),
		ContentID:       int(contentID),
		InteractionType: interactionType,
		Rating:          rating,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	url := fmt.Sprintf("%s/recommendations/interactions", s.recommenderURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		// Não falhar se o motor não estiver disponível - apenas log
		return fmt.Errorf("erro ao notificar motor (não crítico): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro do motor (não crítico) (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

