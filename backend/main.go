package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Handler for the GET request on the "/" route
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

// Handler for the POST request on the "/greet" route
func greetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	type RequestBody struct {
		Name string `json:"name"`
	}

	var reqBody RequestBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	greeting := fmt.Sprintf("Hello, %s!", reqBody.Name)
	response := map[string]string{"greeting": greeting}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	mux := http.NewServeMux()

	// Register the helloHandler function for the "/" route
	mux.HandleFunc("/", helloHandler)

	// Register the greetHandler function for the "/greet" route
	mux.HandleFunc("/greet", greetHandler)

	// Start the server on port 8080
	fmt.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
