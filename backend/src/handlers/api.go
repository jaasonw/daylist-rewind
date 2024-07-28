package handlers

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	"daylist-rewind-backend/src/pocketbase"

	"github.com/go-chi/chi/v5"
)

func GetPlaylistTracksHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the playlist ID from the URL path
	playlistID := chi.URLParam(r, "playlistID")
	if playlistID == "" {
		http.Error(w, "Missing playlist ID", http.StatusBadRequest)
		return
	}

	bearer, err := pocketbase.Authenticate(os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		log.Fatal("Error authenticating: " + err.Error())
	}

	songs, err := pocketbase.GetPlaylistSongs(playlistID, bearer)
	if err != nil {
		slog.Error("Error getting playlist songs: " + err.Error())
		http.Error(w, "Error getting playlist songs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Process the playlist ID and songs, and return a JSON response
	response := map[string]interface{}{
		"playlist_id": playlistID,
		"songs":       songs,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Error encoding JSON response: " + err.Error())
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract userID from the URL path
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusForbidden)
		return
	}
	accessToken := r.URL.Query().Get("access_token")
	if accessToken == "" {
		// 403 forbidden
		http.Error(w, "Missing access token", http.StatusForbidden)
		return
	}

	bearer, err := pocketbase.Authenticate(os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		log.Fatal("Error authenticating: " + err.Error())
	}

	profile, err := pocketbase.GetUserRecord(userID, bearer)
	if err != nil {
		slog.Error("Error getting user profile: " + err.Error())
		http.Error(w, "Error getting user profile: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if profile.AccessToken == accessToken {
		// 200 OK
		if err := json.NewEncoder(w).Encode(profile); err != nil {
			slog.Error("Error encoding JSON response: " + err.Error())
			http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		// 403 forbidden
		http.Error(w, "Invalid access token", http.StatusForbidden)
		return
	}
}

func GetUserPlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract userID from the URL path
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	bearer, err := pocketbase.Authenticate(os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		log.Fatal("Error authenticating: " + err.Error())
	}

	playlists, err := pocketbase.GetUserPlaylists(userID, bearer)
	if err != nil {
		slog.Error("Error getting user playlists: " + err.Error())
		http.Error(w, "Error getting playlists: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Process the playlist ID and songs, and return a JSON response
	if err := json.NewEncoder(w).Encode(playlists); err != nil {
		slog.Error("Error encoding JSON response: " + err.Error())
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetSingleUserPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	// Extract playlistId from the URL path
	playlistId := chi.URLParam(r, "playlistID")
	if playlistId == "" {
		http.Error(w, "Missing playlist ID", http.StatusBadRequest)
		return
	}

	bearer, err := pocketbase.Authenticate(os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		log.Fatal("Error authenticating: " + err.Error())
	}

	playlists, err := pocketbase.GetPlaylistRecord(playlistId, bearer)
	if err != nil {
		slog.Error("Error getting user playlists: " + err.Error())
		http.Error(w, "Error getting playlists: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Process the playlist ID and songs, and return a JSON response
	if err := json.NewEncoder(w).Encode(playlists); err != nil {
		slog.Error("Error encoding JSON response: " + err.Error())
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
	}
}
