package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"daylist-rewind-backend/src/pocketbase"
	"daylist-rewind-backend/src/spotifyutil"
	"daylist-rewind-backend/src/util"

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

// Handler for GetPlaylistRecord
// Path Parameters:
// - playlistID: the db provided ID of the playlist to get
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

// Handler for ExportPlaylist
// Query Parameters:
// - username: the spotify username of the user to export the playlist to
// - access_token: the access token of the user
// Path Parameters:
// - playlistID: the db provided ID of the playlist to export
func ExportPlaylistToSpotifyHandler(w http.ResponseWriter, r *http.Request) {
	// return a json with the id
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":      "error",
		"playlist_id": "",
		"message":     "",
	}

	// get user to export to
	username := r.URL.Query().Get("username")
	if username == "" {
		response["message"] = "Missing username"
		util.LogAndProduceJsonResponse(w, http.StatusForbidden, response)
		return
	}

	accessToken := r.URL.Query().Get("access_token")
	if accessToken == "" {
		response["message"] = "Missing access token"
		util.LogAndProduceJsonResponse(w, http.StatusForbidden, response)
		return
	}

	playlistId := chi.URLParam(r, "playlistID")
	if playlistId == "" {
		response["message"] = "Missing playlist ID"
		util.LogAndProduceJsonResponse(w, http.StatusBadRequest, response)
		return
	}

	bearer, err := pocketbase.Authenticate(os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		log.Fatal("Error authenticating: " + err.Error())
		return
	}

	_, err = pocketbase.ValidateToken(username, accessToken, bearer)
	if err != nil {
		response["message"] = "Invalid Access Token"
		util.LogAndProduceJsonResponse(w, http.StatusForbidden, response)
		return
	}

	// get user record
	user, err := pocketbase.GetUserRecord(username, bearer)
	if err != nil {
		response["message"] = "Error getting user record: " + err.Error()
		util.LogAndProduceJsonResponse(w, http.StatusInternalServerError, response)
		return
	}

	// authenticate with spotify
	client, err := spotifyutil.CreateClient(user.AccessToken, user.RefreshToken, user.Expiry)
	if err != nil {
		response["message"] = "Error authenticating with Spotify"
		util.LogAndProduceJsonResponse(w, http.StatusInternalServerError, response)
		return
	}

	// Export playlist to Spotify
	createdPlaylistId, err := spotifyutil.ExportPlaylist(client, playlistId)
	response["playlist_id"] = createdPlaylistId
	if err != nil {
		response["message"] = "Error exporting playlist to Spotify"
		util.LogAndProduceJsonResponse(w, http.StatusInternalServerError, response)
		return
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		response["message"] = "Error encoding JSON response: (" + fmt.Sprintf("%v", response) + ")" + err.Error()
		util.LogAndProduceJsonResponse(w, http.StatusInternalServerError, response)
	}
}
