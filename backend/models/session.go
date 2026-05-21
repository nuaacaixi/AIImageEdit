package models

import (
	"encoding/json"
	"os"
	"time"
)

type Session struct {
	ID               string    `json:"sessionId"`
	OriginalImageID  string    `json:"originalImageId"`
	OriginalImageURL string    `json:"originalImageUrl"`
	CurrentImageID   string    `json:"currentImageId"`
	CurrentImageURL  string    `json:"currentImageUrl"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type LLMResponse struct {
	Intent          string `json:"intent"`
	ToolName        string `json:"toolName"`
	RewrittenPrompt string `json:"rewrittenPrompt"`
	Reasoning       string `json:"reasoning"`
}

type Turn struct {
	ID             string     `json:"turnId"`
	SessionID      string     `json:"sessionId"`
	UserInput      string     `json:"userInput"`
	LLMResponse    LLMResponse `json:"llmResponse"`
	ToolName       string     `json:"toolName"`
	ToolParamsJSON string     `json:"-"` // stored as JSON string internally
	ResultImageID  string     `json:"resultImageId"`
	ResultImageURL string     `json:"resultImageUrl"`
	InputImageID   string     `json:"inputImageId"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"createdAt"`
}

func (t *Turn) SetToolParams(p map[string]any) error {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	t.ToolParamsJSON = string(b)
	return nil
}

type ContextMessage struct {
	ID        string    `json:"id"`
	SessionID string    `json:"sessionId"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	ImageRef  string    `json:"imageRef,omitempty"`
	TurnID    string    `json:"turnId,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type TurnStatus string

const (
	TurnStatusPending    TurnStatus = "pending"
	TurnStatusProcessing TurnStatus = "processing"
	TurnStatusDone       TurnStatus = "done"
	TurnStatusFailed     TurnStatus = "failed"
)

const (
	ResultPrefix = "edited-"
)

func GenerateResultID(originalID string) string {
	return ResultPrefix + originalID
}

func SaveBytesToFile(data []byte, path string) error {
	return os.WriteFile(path, data, 0o644)
}
