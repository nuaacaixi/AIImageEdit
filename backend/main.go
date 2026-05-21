package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/vibeCoding/AIImageEdit/backend/ai"
	"github.com/vibeCoding/AIImageEdit/backend/handlers"
	"github.com/vibeCoding/AIImageEdit/backend/models"
	"github.com/vibeCoding/AIImageEdit/backend/tools"
)

func main() {
	// Storage for images
	storagePath := filepath.Join(".", "storage", "images")
	if err := os.MkdirAll(storagePath, 0o755); err != nil {
		log.Fatalf("failed to create storage directory: %v", err)
	}

	// Data store for sessions/turns/context
	dataPath := filepath.Join(".", "storage", "data")
	if err := models.InitDB(dataPath); err != nil {
		log.Fatalf("failed to init data store: %v", err)
	}
	store := models.DB()

	// AI image client
	aiClient, err := ai.NewClient()
	if err != nil {
		log.Printf("WARNING: AI image client not initialized: %v", err)
		log.Println("Image editing tools will not work until OPENAI_API_KEY is set.")
	}

	// LLM Gateway
	llmGateway, err := ai.NewLLMGateway()
	if err != nil {
		log.Printf("WARNING: LLM Gateway not initialized: %v", err)
		log.Println("LLM intent parsing will fall back to default edit mode.")
	}

	// Tool registry
	toolReg := tools.NewRegistry()
	if aiClient != nil {
		toolReg.Register(tools.NewEditImageTool(aiClient))
	}
	log.Printf("Tool registry initialized with %d tools: %v", len(toolReg.List()), toolReg.EnabledToolNames())

	// Handlers
	uploadHandler := handlers.NewUploadHandler(storagePath, store)
	chatHandler := handlers.NewChatHandler(store, storagePath, llmGateway, toolReg)
	sessionHandler := handlers.NewSessionHandler(store)

	mux := http.NewServeMux()

	// Upload endpoint
	mux.HandleFunc("/api/upload", uploadHandler.Handle)

	// Chat endpoint (main AI interaction)
	mux.HandleFunc("/api/chat", chatHandler.Handle)

	// Session endpoints
	mux.HandleFunc("/api/session/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > len("/api/session/") {
			suffix := path[len("/api/session/"):]
			if len(suffix) > 6 && suffix[len(suffix)-6:] == "/turns" {
				sessionHandler.HandleTurns(w, r)
				return
			}
			if len(suffix) > 8 && suffix[len(suffix)-8:] == "/context" {
				sessionHandler.HandleContext(w, r)
				return
			}
		}
		sessionHandler.HandleSession(w, r)
	})

	// Static image serving
	mux.Handle("/api/images/", http.StripPrefix("/api/images/", http.FileServer(http.Dir(storagePath))))

	// Legacy edit endpoint (kept for backward compat)
	if aiClient != nil {
		editHandler := handlers.NewEditHandler(storagePath, aiClient)
		mux.HandleFunc("/api/edit", editHandler.Handle)
	}

	log.Println("=== AI Image Editor Backend ===")
	log.Printf("Image storage: %s", storagePath)
	log.Printf("Data store: %s", dataPath)
	log.Printf("AI Image Client: %v", aiClient != nil)
	log.Printf("LLM Gateway: %v (%s)", llmGateway != nil, func() string {
		if llmGateway != nil {
			return llmGateway.Model()
		}
		return "N/A"
	}())
	log.Printf("Tools: %v", toolReg.EnabledToolNames())
	log.Println("Server starting on http://localhost:8080")

	if err := http.ListenAndServe(":8080", loggingMiddleware(corsMiddleware(mux))); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
