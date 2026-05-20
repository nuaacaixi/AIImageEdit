package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/vibeCoding/AIImageEdit/backend/ai"
	"github.com/vibeCoding/AIImageEdit/backend/handlers"
)

func main() {
	storagePath := filepath.Join(".", "storage", "images")
	if err := os.MkdirAll(storagePath, 0o755); err != nil {
		log.Fatalf("failed to create storage directory: %v", err)
	}

	client, err := ai.NewClient()
	if err != nil {
		log.Printf("WARNING: AI client not initialized: %v", err)
		log.Println("The /api/edit endpoint will not work until OPENAI_API_KEY is set.")
	}

	uploadHandler := handlers.NewUploadHandler(storagePath)
	editHandler := handlers.NewEditHandler(storagePath, client)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/upload", uploadHandler.Handle)
	mux.HandleFunc("/api/edit", editHandler.Handle)
	mux.Handle("/api/images/", http.StripPrefix("/api/images/", http.FileServer(http.Dir(storagePath))))

	// If you build the frontend and place the output in frontend/dist,
	// uncomment the static file server below and point to the dist directory.
	// mux.Handle("/", http.FileServer(http.Dir("../frontend/dist")))

	log.Println("starting backend on http://localhost:8080")
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
