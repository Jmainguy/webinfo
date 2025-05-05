package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func main() {
	// Custom handler to set MIME type for .wasm files
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".wasm") {
			w.Header().Set("Content-Type", "application/wasm")
		}
		http.ServeFile(w, r, filepath.Join("/app", r.URL.Path))
	})

	// Start the HTTP server on port 8080
	log.Println("Serving on http://0.0.0.0:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
