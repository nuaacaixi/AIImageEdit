package models

import "fmt"

const MaxContextMessages = 20

type LLMMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

func BuildLLMContext(store *Store, sessionID string, systemPrompt string, userInput string, baseImageURL string) ([]LLMMessage, error) {
	msgs, err := store.GetContextMessages(sessionID, MaxContextMessages)
	if err != nil {
		return nil, fmt.Errorf("get context messages: %w", err)
	}

	messages := []LLMMessage{}

	messages = append(messages, LLMMessage{
		Role:    "system",
		Content: systemPrompt,
	})

	for _, m := range msgs {
		msg := LLMMessage{Role: m.Role}
		if m.ImageRef != "" {
			msg.Content = []ContentPart{
				{Type: "text", Text: m.Content},
				{Type: "image_url", ImageURL: &ImageURL{URL: m.ImageRef, Detail: "auto"}},
			}
		} else {
			msg.Content = m.Content
		}
		messages = append(messages, msg)
	}

	if len(msgs) == 0 && baseImageURL != "" {
		messages = append(messages, LLMMessage{
			Role: "user",
			Content: []ContentPart{
				{Type: "text", Text: userInput},
				{Type: "image_url", ImageURL: &ImageURL{URL: baseImageURL, Detail: "auto"}},
			},
		})
	} else {
		messages = append(messages, LLMMessage{
			Role:    "user",
			Content: userInput,
		})
	}

	return messages, nil
}
