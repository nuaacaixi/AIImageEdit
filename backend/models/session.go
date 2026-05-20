package models

import (
    "fmt"
    "os"
)

type EditResult struct {
    ImageID   string `json:"imageId"`
    Prompt    string `json:"prompt"`
    ResultURL string `json:"resultUrl"`
}

type SessionRecord struct {
    SessionID string       `json:"sessionId"`
    ImageID   string       `json:"imageId"`
    Results   []EditResult `json:"results"`
}

func GenerateResultID(originalID string) string {
    return fmt.Sprintf("edited-%s", originalID)
}

func SaveBytesToFile(data []byte, path string) error {
    return os.WriteFile(path, data, 0o644)
}
