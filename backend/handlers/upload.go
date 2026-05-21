package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/vibeCoding/AIImageEdit/backend/models"
)

type UploadHandler struct {
	storagePath string
	store       *models.Store
}

type UploadResponse struct {
	ImageID     string `json:"imageId"`
	OriginalURL string `json:"originalUrl"`
	SessionID   string `json:"sessionId"`
}

func NewUploadHandler(storagePath string, store *models.Store) *UploadHandler {
	return &UploadHandler{storagePath: storagePath, store: store}
}

func (h *UploadHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "image file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageID := generateImageID(header)
	destination := filepath.Join(h.storagePath, imageID)
	if err := saveFile(file, destination); err != nil {
		http.Error(w, "failed to save image", http.StatusInternalServerError)
		return
	}

	originalURL := fmt.Sprintf("/api/images/%s", imageID)

	// Create a new session for this upload
	sessionID := generateUUID()
	session := &models.Session{
		ID:               sessionID,
		OriginalImageID:  imageID,
		OriginalImageURL: originalURL,
		CurrentImageID:   imageID,
		CurrentImageURL:  originalURL,
	}

	if err := h.store.CreateSession(session); err != nil {
		// Session creation failed but image was uploaded,
		// still return the image info
		response := UploadResponse{
			ImageID:     imageID,
			OriginalURL: originalURL,
			SessionID:   sessionID,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Add initial context message
	h.store.AddContextMessage(&models.ContextMessage{
		ID:        generateUUID(),
		SessionID: sessionID,
		Role:      "assistant",
		Content:   "图片已上传成功！我是您的 AI 修图助手，请告诉我您想如何调整这张图片。",
		ImageRef:  originalURL,
	})

	response := UploadResponse{
		ImageID:     imageID,
		OriginalURL: originalURL,
		SessionID:   sessionID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateImageID(header *multipart.FileHeader) string {
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".png"
	}
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
}

func saveFile(file multipart.File, destination string) error {
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	return err
}
