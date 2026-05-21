package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/vibeCoding/AIImageEdit/backend/models"
)

type SessionHandler struct {
	store *models.Store
}

func NewSessionHandler(store *models.Store) *SessionHandler {
	return &SessionHandler{store: store}
}

func (h *SessionHandler) HandleSession(w http.ResponseWriter, r *http.Request) {
	// Extract session ID: /api/session/{id}
	id := strings.TrimPrefix(r.URL.Path, "/api/session/")
	id = strings.TrimSuffix(id, "/turns")
	id = strings.TrimSuffix(id, "/context")
	if id == "" || id == r.URL.Path {
		http.Error(w, "session id is required", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sess, err := h.store.GetSession(id)
	if err != nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sess)
}

func (h *SessionHandler) HandleTurns(w http.ResponseWriter, r *http.Request) {
	// Extract session ID: /api/session/{id}/turns
	id := strings.TrimPrefix(r.URL.Path, "/api/session/")
	id = strings.TrimSuffix(id, "/turns")
	if id == "" || id == r.URL.Path {
		http.Error(w, "session id is required", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	turns, err := h.store.GetTurnsBySession(id)
	if err != nil {
		http.Error(w, "failed to get turns: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(turns)
}

func (h *SessionHandler) HandleContext(w http.ResponseWriter, r *http.Request) {
	// Extract session ID: /api/session/{id}/context
	id := strings.TrimPrefix(r.URL.Path, "/api/session/")
	id = strings.TrimSuffix(id, "/context")
	if id == "" || id == r.URL.Path {
		http.Error(w, "session id is required", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	msgs, err := h.store.GetContextMessages(id, models.MaxContextMessages)
	if err != nil {
		http.Error(w, "failed to get context: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}
