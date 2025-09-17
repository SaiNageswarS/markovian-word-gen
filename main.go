package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

//go:embed views/*
var viewsFS embed.FS

// GenerateWordsRequest represents the request structure for generating words
type GenerateWordsRequest struct {
	Words     []string `json:"words"`
	Count     int      `json:"count,omitempty"`     // number of words to generate (default: 10)
	MaxLength int      `json:"maxLength,omitempty"` // maximum length of generated words (default: 10)
}

// GenerateWordsResponse represents the response structure
type GenerateWordsResponse struct {
	GeneratedWords []string `json:"generatedWords"`
	Count          int      `json:"count"`
	MaxLength      int      `json:"maxLength"`
}

// generateWordsHandler handles the /generateWords endpoint
func generateWordsHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST requests
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req GenerateWordsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if len(req.Words) == 0 {
		http.Error(w, "Words list cannot be empty", http.StatusBadRequest)
		return
	}

	// Set defaults
	if req.Count <= 0 {
		req.Count = 10
	}
	if req.MaxLength <= 0 {
		req.MaxLength = 10
	}

	// Create word generator and build from input words
	generator := NewWordGenerator()
	generator.Build(req.Words)

	// Generate words
	generatedWords := make([]string, req.Count)
	for i := 0; i < req.Count; i++ {
		generatedWords[i] = generator.Generate(req.MaxLength)
	}

	// Create response
	response := GenerateWordsResponse{
		GeneratedWords: generatedWords,
		Count:          req.Count,
		MaxLength:      req.MaxLength,
	}

	// Send response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// healthHandler provides a simple health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// indexHandler serves the main HTML page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Read the embedded template file
	templateData, err := viewsFS.ReadFile("views/index.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.New("index").Parse(string(templateData)))
	tmpl.Execute(w, nil)
}

func main() {
	// Set up routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/generateWords", generateWordsHandler)
	http.HandleFunc("/health", healthHandler)

	// Start server
	port := ":8080"
	fmt.Printf("Starting HTTP server on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /             - Web interface for word generation")
	fmt.Println("  POST /generateWords - Generate words based on input vocabulary")
	fmt.Println("  GET  /health       - Health check endpoint")

	log.Fatal(http.ListenAndServe(port, nil))
}
