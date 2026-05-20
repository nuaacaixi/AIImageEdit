package ai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Client struct {
	apiKey string
	apiURL string
	model  string
}

type openAIResponse struct {
	Data []struct {
		B64JSON string `json:"b64_json"`
	} `json:"data"`
}

func NewClient() (*Client, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	apiURL := os.Getenv("OPENAI_API_BASE")
	if apiURL == "" {
		apiURL = "https://api.openai.com/v1"
	}

	return &Client{
		apiKey: apiKey,
		apiURL: apiURL,
		model:  "gpt-image-1",
	}, nil
}

func (c *Client) EditImage(imagePath, prompt string) ([]byte, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("open image: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("model", c.model); err != nil {
		return nil, err
	}
	if err := writer.WriteField("prompt", prompt); err != nil {
		return nil, err
	}

	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, c.apiURL+"/images/edits", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("third-party API error: %s", string(data))
	}

	var parsed openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}
	if len(parsed.Data) == 0 || parsed.Data[0].B64JSON == "" {
		return nil, errors.New("missing image data from third-party API")
	}

	return base64.StdEncoding.DecodeString(parsed.Data[0].B64JSON)
}
