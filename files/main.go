package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const baseDir = "./data"

func main() {
	http.HandleFunc("/", handleRequest)

	log.Println("Serving files and accepting PUT requests at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPut:
		handlePut(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGet serves files from the `baseDir` directory.
func handleGet(w http.ResponseWriter, r *http.Request) {
	// Clean and convert the URL path into an OS-friendly path
	path := filepath.Join(baseDir, filepath.FromSlash(strings.TrimPrefix(r.URL.Path, "/")))

	// Use http.ServeFile for convenience
	http.ServeFile(w, r, path)
}

// handlePut writes the request body to a file under the `baseDir` directory.
func handlePut(w http.ResponseWriter, r *http.Request) {
	// Clean and convert the URL path
	path := filepath.Join(baseDir, filepath.FromSlash(strings.TrimPrefix(r.URL.Path, "/")))

	// Create the necessary directories
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		http.Error(w, "Could not create directories", http.StatusInternalServerError)
		return
	}

	// Create or overwrite the file
	file, err := os.Create(path)
	if err != nil {
		http.Error(w, "Could not create file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Copy the request body into the file
	if _, err := io.Copy(file, r.Body); err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}

	// Indicate that the resource was created or overwritten
	w.WriteHeader(http.StatusCreated)
}
