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
)

type UploadHandler struct {
    storagePath string
}

type UploadResponse struct {
    ImageID     string `json:"imageId"`
    OriginalURL string `json:"originalUrl"`
}

func NewUploadHandler(storagePath string) *UploadHandler {
    return &UploadHandler{storagePath: storagePath}
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

    imageID := generateID(header)
    destination := filepath.Join(h.storagePath, imageID)
    if err := saveFile(file, destination); err != nil {
        http.Error(w, "failed to save image", http.StatusInternalServerError)
        return
    }

    response := UploadResponse{
        ImageID:     imageID,
        OriginalURL: fmt.Sprintf("/api/images/%s", imageID),
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func generateID(header *multipart.FileHeader) string {
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
