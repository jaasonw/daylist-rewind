package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// PostRequest sends a POST request with a given body (any struct) and returns the response body as a byte slice.
func PostRequest(url string, body interface{}, headers map[string]string) ([]byte, error) {
	// Marshal the body to JSON
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Add additional headers to the request if provided
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Send the request using an http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return respBody, fmt.Errorf("received non-200 response status: %s", resp.Status)
	}

	return respBody, nil
}

// GetRequest sends a GET request and returns the response body as a byte slice.
func GetRequest(url string, query map[string]string, headers map[string]string) ([]byte, error) {
	// Create the request URL with query parameters
	if len(query) > 0 {
		url += "?"
		for key, value := range query {
			url += key + "=" + value + "&"
		}
		url = url[:len(url)-1]
	}

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %v", err)
	}

	// Add headers to the request
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Send the request using an http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %s", resp.Status)
	}

	return respBody, nil
}

func PatchRequest(url string, body interface{}, headers map[string]string) ([]byte, error) {
	// Marshal the body to JSON
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create a new PATCH request
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create PATCH request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Add additional headers to the request if provided
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Send the request using an http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make PATCH request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %s", resp.Status)
	}

	return respBody, nil
}
