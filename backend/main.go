// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )

// // Handler for the GET request on the "/" route
// func helloHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Hello, World!")
// }

// // Handler for the POST request on the "/greet" route
// func greetHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	type RequestBody struct {
// 		Name string `json:"name"`
// 	}

// 	var reqBody RequestBody
// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&reqBody); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	greeting := fmt.Sprintf("Hello, %s!", reqBody.Name)
// 	response := map[string]string{"greeting": greeting}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// func main() {
// 	mux := http.NewServeMux()

// 	// Register the helloHandler function for the "/" route
// 	mux.HandleFunc("/", helloHandler)

// 	// Register the greetHandler function for the "/greet" route
// 	mux.HandleFunc("/greet", greetHandler)

// 	// Start the server on port 8080
// 	fmt.Println("Starting server on port 8080...")
// 	if err := http.ListenAndServe(":8080", mux); err != nil {
// 		fmt.Println("Error starting server:", err)
// 	}
// }
// Command profile gets the public profile information about a Spotify user.

// package main

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"os"

// 	spotifyauth "github.com/zmb3/spotify/v2/auth"

// 	"github.com/joho/godotenv"
// 	"github.com/zmb3/spotify/v2"
// 	"golang.org/x/oauth2/clientcredentials"
// )

// var userID = flag.String("user", "", "the Spotify user ID to look up")

// func main() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	flag.Parse()

// 	ctx := context.Background()

// 	if *userID == "" {
// 		fmt.Fprintf(os.Stderr, "Error: missing user ID\n")
// 		flag.Usage()
// 		return
// 	}

// 	config := &clientcredentials.Config{
// 		ClientID:     os.Getenv("SPOTIFY_ID"),
// 		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
// 		TokenURL:     spotifyauth.TokenURL,
// 	}
// 	// token, err := config.Token(context.Background())
// 	if err != nil {
// 		log.Fatalf("couldn't get token: %v", err)
// 	}

// 	httpClient := spotifyauth.New().Client(ctx, token)
// 	// client := spotify.New(httpClient)
// 	// client := spotify.New(auth.Client(r.Context(), "BQAtrGpXgFqQdJVdOQksQJuibARbfGvOwO8W6RqLEPOlu1Z-iCt4KBm0jMN4aaOCvHRiA2Zuehq7VGgIxoRlFJBmkWg4b4UBNBpTBXGZEFXMF-JqRvtkt1u3nLsp4b6KsBirk6NQS3CdpqQdqFAchLOIf_QnEyCplhbNrQGkJcC7U531rHIpQMxoZsEHJCRXfi_dXwiTbroao64ARPfFXzZHTXE-HbIb9uB_NN9G"))

// 	auth  := spotifyauth.New(spotifyauth.WithRedirectURL(), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
// 	token := "BQAtrGpXgFqQdJVdOQksQJuibARbfGvOwO8W6RqLEPOlu1Z-iCt4KBm0jMN4aaOCvHRiA2Zuehq7VGgIxoRlFJBmkWg4b4UBNBpTBXGZEFXMF-JqRvtkt1u3nLsp4b6KsBirk6NQS3CdpqQdqFAchLOIf_QnEyCplhbNrQGkJcC7U531rHIpQMxoZsEHJCRXfi_dXwiTbroao64ARPfFXzZHTXE-HbIb9uB_NN9G"
// 	client := spotify.New(auth.Client(context.Background(), &token))
// 	// client := spotify.Authenticator{}.NewClient("client := spotify.Authenticator{}.NewClient(accessToken)")
// 	user, err := client.GetUsersPublicProfile(ctx, spotify.ID(*userID))
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err.Error())
// 		return
// 	}

// 	fmt.Println("User ID:", user.ID)
// 	fmt.Println("Display name:", user.DisplayName)
// 	fmt.Println("Spotify URI:", string(user.URI))
// 	fmt.Println("Endpoint:", user.Endpoint)
// 	fmt.Println("Followers:", user.Followers.Count)
// }

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
)

func pb_authenticate(identity, password string) (string, error) {

	type AuthRequest struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	type AuthResponse struct {
		Token string `json:"token"`
	}

	url := "http://localhost:8090/api/admins/auth-with-password"

	// Create the request body
	authReq := AuthRequest{
		Identity: identity,
		Password: password,
	}

	reqBody, err := json.Marshal(authReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create the POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response status: %s", resp.Status)
	}

	// Unmarshal the response body into AuthResponse
	var authResp AuthResponse
	err = json.Unmarshal(respBody, &authResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return authResp.Token, nil
}

type Item struct {
	ID              string `json:"id"`
	CollectionId    string `json:"collectionId"`
	CollectionName  string `json:"collectionName"`
	Username        string `json:"username"`
	Verified        bool   `json:"verified"`
	EmailVisibility bool   `json:"emailVisibility"`
	Email           string `json:"email"`
	Created         string `json:"created"`
	Updated         string `json:"updated"`
	SpotifyUsername string `json:"spotify_username"`
	SpotifyEmail    string `json:"spotify_email"`
	SpotifyID       string `json:"spotify_id"`
	AccessToken     string `json:"accessToken"`
	RefreshToken    string `json:"refreshToken"`
	Expiry          string `json:"expiry"`
	DisplayName     string `json:"display_name"`
	AvatarURL       string `json:"avatar_url"`
	ErrorCount      int    `json:"error_count"`
	LastError       string `json:"last_error"`
	Active          bool   `json:"active"`
}

type RecordsResponse struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"perPage"`
	TotalPages int    `json:"totalPages"`
	TotalItems int    `json:"totalItems"`
	Items      []Item `json:"items"`
}

func getAllRecords(token string) ([]Item, error) {

	var allItems []Item
	url := "http://localhost:8090/api/collections/users/records"
	page := 1

	for {
		// Create the request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		// Set the authorization header
		req.Header.Set("Authorization", "Bearer "+token)

		// Add the query parameters
		q := req.URL.Query()
		q.Add("perPage", "500")
		q.Add("page", fmt.Sprintf("%d", page))
		req.URL.RawQuery = q.Encode()

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make GET request: %v", err)
		}
		defer resp.Body.Close()

		// Read the response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("received non-200 response status: %s", resp.Status)
		}

		// Unmarshal the response body into RecordsResponse
		var recordsResp RecordsResponse
		err = json.Unmarshal(respBody, &recordsResp)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
		}

		// Append items to the list
		allItems = append(allItems, recordsResp.Items...)

		// Break if we have fetched all pages
		if page >= recordsResp.TotalPages {
			break
		}

		// Go to the next page
		page++
	}

	return allItems, nil
}

func createClient(access_token string, refresh_token string, expiry string) *spotify.Client {

	layout := "2006-01-02 15:04:05.000Z"
	expiryTime, err := time.Parse(layout, expiry)
	if err != nil {
		log.Fatal("Error parsing time" + err.Error())
	}

	token := &oauth2.Token{
		AccessToken:  access_token,
		TokenType:    "Bearer",
		RefreshToken: refresh_token,
		Expiry:       expiryTime,
	}
	token, err = auth.RefreshToken(context.Background(), token)
	if err != nil {
		log.Fatalf("Failed to refresh token: %v", err)
	}
	fmt.Println(token)
	client := spotify.New(auth.Client(context.Background(), token))
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatalf("Failed to get user info: %v", err)
	}

	// Print the user info
	fmt.Printf("User ID: %s\n", user.ID)
	return client
}

func main() {
	fmt.Println(uuid.New().String())
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	access_token := ""
	refresh_token := ""
	expiry := ""

	client := createClient(access_token, refresh_token, expiry)

	// This playlist ID is constant it is a spotify owned playlist that is unique to whoever is logged in
	daylist, err := client.GetPlaylist(context.Background(), spotify.ID("37i9dQZF1EP6YuccBxUcC1"))
	if err != nil {
		log.Fatal("Error getting playlist" + err.Error())
	}
	fmt.Println(daylist.Name)
	fmt.Println(daylist.Description)
	fmt.Println(daylist.Owner.DisplayName)
	for _, track := range daylist.Tracks.Tracks {
		fmt.Println(
			track.Track.Name,
			track.Track.Artists[0].Name,
			track.Track.Album.Name,
			track.Track.Album.Images[0].URL,
		)
	}
	bearer, err := pb_authenticate(os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		log.Fatal("Error authenticating" + err.Error())
	}
	fmt.Println(bearer)

	items, err := getAllRecords(bearer)
	if err != nil {
		log.Fatal("Error getting records" + err.Error())
	}
	for _, item := range items {
		fmt.Println(item.ID)
	}
}
