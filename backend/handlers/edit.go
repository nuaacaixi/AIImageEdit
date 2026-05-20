package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/vibeCoding/AIImageEdit/backend/ai"
	"github.com/vibeCoding/AIImageEdit/backend/models"
)

type EditHandler struct {
	storagePath string
	client      *ai.Client
}

type EditResponse struct {
	Status    string `json:"status"`
	ResultURL string `json:"resultUrl"`
	Message   string `json:"message"`
}

type editRequest struct {
	ImageID string `json:"imageId"`
	Prompt  string `json:"prompt"`
}

func NewEditHandler(storagePath string, client *ai.Client) *EditHandler {
	return &EditHandler{storagePath: storagePath, client: client}
}

func (h *EditHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req editRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.ImageID == "" || req.Prompt == "" {
		http.Error(w, "imageId and prompt are required", http.StatusBadRequest)
		return
	}

	if h.client == nil {
		http.Error(w, "AI service is not configured: OPENAI_API_KEY is missing", http.StatusServiceUnavailable)
		return
	}

	originalPath := filepath.Join(h.storagePath, req.ImageID)
	resultBytes, err := h.client.EditImage(originalPath, req.Prompt)
	if err != nil {
		http.Error(w, "failed to edit image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resultID := models.GenerateResultID(req.ImageID)
	resultPath := filepath.Join(h.storagePath, resultID)
	if err := models.SaveBytesToFile(resultBytes, resultPath); err != nil {
		http.Error(w, "failed to save edited image", http.StatusInternalServerError)
		return
	}

	response := EditResponse{
		Status:    "success",
		ResultURL: "/api/images/" + resultID,
		Message:   "Image edited successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
