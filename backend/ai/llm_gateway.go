package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/vibeCoding/AIImageEdit/backend/models"
)

type LLMGateway struct {
	apiKey  string
	apiURL  string
	model   string
	client  *http.Client
}

type LLMResponsePayload struct {
	Intent          string `json:"intent"`
	ToolName        string `json:"toolName"`
	RewrittenPrompt string `json:"rewrittenPrompt"`
	Reasoning       string `json:"reasoning"`
}

type chatRequest struct {
	Model    string              `json:"model"`
	Messages []models.LLMMessage `json:"messages"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewLLMGateway() (*LLMGateway, error) {
	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("LLM_API_KEY or OPENAI_API_KEY is required for LLM Gateway")
	}

	apiURL := os.Getenv("LLM_API_BASE")
	if apiURL == "" {
		apiURL = os.Getenv("OPENAI_API_BASE")
	}
	if apiURL == "" {
		apiURL = "https://api.openai.com/v1"
	}

	model := os.Getenv("LLM_MODEL")
	if model == "" {
		model = "gpt-4o"
	}

	return &LLMGateway{
		apiKey: apiKey,
		apiURL: apiURL,
		model:  model,
		client: &http.Client{},
	}, nil
}

func (g *LLMGateway) ParseIntent(messages []models.LLMMessage) (*LLMResponsePayload, error) {
	body := chatRequest{
		Model:    g.model,
		Messages: messages,
	}
	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, g.apiURL+"/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("llm request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("llm api error (%d): %s", resp.StatusCode, string(data))
	}

	var llmResp chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&llmResp); err != nil {
		return nil, fmt.Errorf("decode llm response: %w", err)
	}
	if len(llmResp.Choices) == 0 {
		return nil, fmt.Errorf("llm returned empty response")
	}

	content := llmResp.Choices[0].Message.Content
	payload, err := parseStructuredResponse(content)
	if err != nil {
		return nil, fmt.Errorf("parse llm output: %w", err)
	}
	return payload, nil
}

func parseStructuredResponse(content string) (*LLMResponsePayload, error) {
	// Extract JSON from potentially markdown-wrapped response
	jsonStr := content
	// Try to find JSON block
	if idx := findJSONStart(content); idx >= 0 {
		end := findJSONEnd(content, idx)
		if end > idx {
			jsonStr = content[idx:end]
		}
	}

	var payload LLMResponsePayload
	if err := json.Unmarshal([]byte(jsonStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid JSON from LLM: %s", content)
	}
	if payload.ToolName == "" {
		return nil, fmt.Errorf("LLM response missing toolName: %s", content)
	}
	return &payload, nil
}

func findJSONStart(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == '{' {
			return i
		}
	}
	return -1
}

func findJSONEnd(s string, start int) int {
	depth := 0
	for i := start; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}
	return len(s)
}

func (g *LLMGateway) Model() string {
	return g.model
}
