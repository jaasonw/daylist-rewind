package pocketbase

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"daylist-rewind-backend/httputil"
)

type AuthRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type UserRecord struct {
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
	Page       int          `json:"page"`
	PerPage    int          `json:"perPage"`
	TotalPages int          `json:"totalPages"`
	TotalItems int          `json:"totalItems"`
	Items      []UserRecord `json:"items"`
}

type Song struct {
	ID             string `json:"id"`
	CollectionId   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	Created        string `json:"created"`
	Updated        string `json:"updated"`
	SongID         string `json:"song_id"`
	Name           string `json:"name"`
	Artist         string `json:"artist"`
	ArtistURL      string `json:"artist_url"`
	Album          string `json:"album"`
	AlbumURL       string `json:"album_url"`
	AlbumCover     string `json:"album_cover"`
	PreviewURL     string `json:"preview_url"`
	Duration       int    `json:"duration"`
}

type Playlist struct {
	ID             string `json:"id"`
	CollectionId   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	Created        string `json:"created"`
	Updated        string `json:"updated"`
	Owner          string `json:"owner"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Hash           string `json:"hash"`
	PlaylistId     string `json:"playlist_id"`
	Data           string `json:"data"`
}

type InsertSongErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		SongID struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"song_id"`
	} `json:"data"`
}

type SongPlaylistLink struct {
	ID             string `json:"id"`
	CollectionId   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	Created        string `json:"created"`
	Updated        string `json:"updated"`
	PlaylistID     string `json:"playlist_id"`
	SongID         string `json:"song_id"`
}

type GetSongResponse struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"perPage"`
	TotalPages int    `json:"totalPages"`
	TotalItems int    `json:"totalItems"`
	Items      []Song `json:"items"`
}

type GetPlaylistResponse struct {
	Page       int        `json:"page"`
	PerPage    int        `json:"perPage"`
	TotalPages int        `json:"totalPages"`
	TotalItems int        `json:"totalItems"`
	Items      []Playlist `json:"items"`
}

type SongPlaylistLinkResponse struct {
	Page       int                `json:"page"`
	PerPage    int                `json:"perPage"`
	TotalPages int                `json:"totalPages"`
	TotalItems int                `json:"totalItems"`
	Items      []SongPlaylistLink `json:"items"`
}

func Authenticate(identity, password string) (string, error) {
	url := "http://localhost:8090/api/admins/auth-with-password"

	// Create the request body
	authReq := AuthRequest{
		Identity: identity,
		Password: password,
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	resp, err := httputil.PostRequest(url, authReq, headers)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate: %v", err)
	}

	// Unmarshal the response body into AuthResponse
	var authResp AuthResponse
	err = json.Unmarshal(resp, &authResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return authResp.Token, nil
}

func GetAllUsers(token string) ([]UserRecord, error) {
	url := "http://localhost:8090/api/collections/users/records"
	page := 1

	var allItems []UserRecord

	for {
		query := map[string]string{
			"perPage": "500",
		}

		headers := map[string]string{
			"Authorization": token,
		}
		response, err := httputil.GetRequest(url, query, headers)
		if err != nil {
			return nil, fmt.Errorf("failed to get users: %v", err)
		}

		// Unmarshal the response body into RecordsResponse
		var recordsResp RecordsResponse
		err = json.Unmarshal(response, &recordsResp)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal users response body: %v", err)
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

// InsertSong inserts a song into the database and returns the song ID.
func InsertSong(song Song, token string) (string, error) {
	url := "http://localhost:8090/api/collections/songs/records"
	headers := map[string]string{
		"Authorization": token,
	}
	response, error := httputil.PostRequest(url, song, headers)
	if error != nil {
		var errorResponse InsertSongErrorResponse
		unmarshalError := json.Unmarshal(response, &errorResponse)
		if unmarshalError != nil {
			return "", fmt.Errorf("failed to unmarshal error response: %v", response)
		}
		return errorResponse.Data.SongID.Message, fmt.Errorf("failed to insert song: %v", error)
	}

	songResponse := Song{}
	err := json.Unmarshal(response, &songResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal song response: %v", err)
	}
	return songResponse.ID, nil
}

// UpdateUser updates a user in the database and returns the updated record
func UpdateUser(user UserRecord, token string) (string, error) {
	url := "http://localhost:8090/api/collections/users/records/" + user.ID
	header := map[string]string{
		"Authorization": token,
	}
	response, error := httputil.PatchRequest(url, user, header)
	if error != nil {
		return "", fmt.Errorf("failed to update user: %v", error)
	}
	return string(response), nil
}

// CreatePlaylist creates a playlist in the database and returns the playlist ID.
func CreatePlaylist(playlist Playlist, token string) (string, error) {
	url := "http://localhost:8090/api/collections/playlists/records"
	headers := map[string]string{
		"Authorization": token,
	}
	response, error := httputil.PostRequest(url, playlist, headers)
	if error != nil {
		return "", fmt.Errorf("failed to create playlist: %v", error)
	}

	playlistResponse := Playlist{}
	err := json.Unmarshal(response, &playlistResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal playlist response: %v", err)
	}
	return playlistResponse.ID, nil
}

// AddSongToPlaylist adds a song to a playlist in the database.
func AddSongToPlaylist(playlistID string, songID string, token string) (string, error) {
	url := "http://localhost:8090/api/collections/song_playlist_link/records"
	headers := map[string]string{
		"Authorization": token,
	}
	body := SongPlaylistLink{
		PlaylistID: playlistID,
		SongID:     songID,
	}
	response, error := httputil.PostRequest(url, body, headers)
	if error != nil {
		return "", fmt.Errorf("failed to insert playlist song link: %v", error)
	}
	return string(response), nil
}

func GetSongBySongId(songID string, token string) (Song, error) {
	url := "http://localhost:8090/api/collections/songs/records"
	headers := map[string]string{
		"Authorization": token,
	}

	query := map[string]string{
		"perPage": "1",
		"filter":  "(song_id='" + songID + "')",
	}
	response, err := httputil.GetRequest(url, query, headers)
	if err != nil {
		return Song{}, fmt.Errorf("failed to get song: %v", err)
	}

	songResponse := GetSongResponse{}
	slog.Debug("response", "response", string(response))
	err = json.Unmarshal(response, &songResponse)
	if err != nil {
		return Song{}, fmt.Errorf("failed to unmarshal SongResponse: %v", err)
	}
	song := songResponse.Items[0]

	return song, nil
}

func CheckPlaylistExists(hash string, token string) (bool, string, error) {
	url := "http://localhost:8090/api/collections/playlists/records"
	headers := map[string]string{
		"Authorization": token,
	}

	query := map[string]string{
		"perPage": "1",
		"filter":  "(hash='" + hash + "')",
	}
	response, err := httputil.GetRequest(url, query, headers)
	if err != nil {
		return false, "", fmt.Errorf("failed to get playlist: %v", err)
	}

	playlistResponse := RecordsResponse{}
	err = json.Unmarshal(response, &playlistResponse)
	if err != nil {
		return false, "", fmt.Errorf("failed to unmarshal playlist response: %v", err)
	}

	if playlistResponse.TotalItems == 0 {
		return false, "", nil
	}

	return true, playlistResponse.Items[0].ID, nil
}

func GetUserPlaylists(userID string, token string) ([]Playlist, error) {
	url := "http://localhost:8090/api/collections/playlists/records"
	headers := map[string]string{
		"Authorization": token,
	}

	query := map[string]string{
		"perPage": "1000",
		"filter":  "(owner='" + userID + "')",
	}
	response, err := httputil.GetRequest(url, query, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlists: %v", err)
	}

	playlistResponse := GetPlaylistResponse{}
	err = json.Unmarshal(response, &playlistResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal playlist response: %v", err)
	}

	return playlistResponse.Items, nil
}

func GetPlaylistSongs(playlistID string, token string) ([]Song, error) {
	url := "http://localhost:8090/api/collections/song_playlist_link/records"
	headers := map[string]string{
		"Authorization": token,
	}

	query := map[string]string{
		"perPage": "1000",
		"filter":  "(playlist_id='" + playlistID + "')",
	}
	response, err := httputil.GetRequest(url, query, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist songs: %v", err)
	}

	songPlaylistResponse := SongPlaylistLinkResponse{}
	err = json.Unmarshal(response, &songPlaylistResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal song playlist response: %v", err)
	}

	var songs []Song
	for _, songPlaylistLink := range songPlaylistResponse.Items {
		song, err := GetSongBySongId(songPlaylistLink.SongID, token)
		if err != nil {
			return nil, fmt.Errorf("failed to get song: %v", err)
		}
		songs = append(songs, song)
	}

	return songs, nil
}
