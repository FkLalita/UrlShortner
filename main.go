package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// url is a map to store the original URLs associated with short keys.
var url = make(map[string]string)

func main() {
	fmt.Println("Starting main function")
	// Register handlers for different routes.
	http.HandleFunc("/", handleForm)           // Handle the root route, displaying the form.
	http.HandleFunc("/shortn", handleShorten)  // Handle form submission to shorten URLs.
	http.HandleFunc("/short/", handleRedirect) // Handle redirection based on short keys.

	fmt.Println("Starting Server")
	// Start the HTTP server on port 8080.
	http.ListenAndServe(":8080", nil)
}

// handleForm handles the root route ("/") to display the URL shortening form.
func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// If the form is submitted, redirect to "/shortn".
		http.Redirect(w, r, "/shortn", http.StatusSeeOther)
	}

	tmpl, err := template.ParseFiles("static/index.html")
	tmpl.Execute(w, url)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Failed to Parse Error")
	}
}

// handleShorten handles form submissions to shorten URLs.
func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Operation", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "No Url Found", http.StatusNotFound)
		return
	}

	// Generate a unique shortKey for the URL.
	shortKey := generateShortKey()
	log.Printf("Generated ShortKey: %s\n", shortKey)

	// Store the original URL with the shortKey in the map.
	url[shortKey] = originalURL

	data := struct {
		OriginalURL string
		ShortenURL  string
	}{
		OriginalURL: originalURL,
		ShortenURL:  fmt.Sprintf("/short/%s", shortKey), // Construct the shortened URL.
	}

	tmpl, err := template.ParseFiles("static/short.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Failed To Parse Template")
		return
	}
	tmpl.Execute(w, data)
}

// handleRedirect handles redirection based on short keys.
func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := strings.TrimPrefix(r.URL.Path, "/short/")
	if shortKey == "" {
		log.Println("ShortKey Not Found")
		return
	}

	originalURL, found := url[shortKey]
	if !found {
		http.Error(w, "ShortKey Not Found", http.StatusNotFound)
		log.Println(shortKey)
		return
	}

	// Redirect the user to the original URL.
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

// generateShortKey generates a unique short key for URLs.
func generateShortKey() string {
	const char = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
	const keyLength = 8
	rand.Seed(time.Now().UnixNano())

	for {
		shortKey := make([]byte, keyLength)
		for i := range shortKey {
			shortKey[i] = char[rand.Intn(len(char))]
		}

		if _, exists := url[string(shortKey)]; !exists {
			return string(shortKey)
		}
	}
}
