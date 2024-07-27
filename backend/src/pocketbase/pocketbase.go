package pocketbase

import (
	"daylist-rewind-backend/src/global"
	"daylist-rewind-backend/src/http"
	"daylist-rewind-backend/src/util"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
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
	AccessToken     string `json:"accessToken"`
	RefreshToken    string `json:"refreshToken"`
	Expiry          string `json:"expiry"`
	DisplayName     string `json:"display_name"`
	AvatarURL       string `json:"avatar_url"`
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
	Expand         struct {
		SongID Song `json:"song_id"`
	} `json:"expand"`
	SongID string `json:"song_id"`
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

type PlaylistSongsResponse struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalPages int `json:"totalPages"`
	TotalItems int `json:"totalItems"`

	Items []SongPlaylistLink `json:"items"`
}

func Authenticate(identity, password string) (string, error) {
	if global.AdminToken != "" && util.ValidateJWT() {
		return global.AdminToken, nil
	}

	url := os.Getenv("POCKETBASE_URL") + "/api/admins/auth-with-password"

	// Create the request body
	authReq := AuthRequest{
		Identity: identity,
		Password: password,
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	resp, err := http.PostRequest(url, authReq, headers)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate: %v", err)
	}

	// Unmarshal the response body into AuthResponse
	var authResp AuthResponse
	err = json.Unmarshal(resp, &authResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	slog.Info("authResp", "authResp", authResp.Token)
	global.AdminToken = authResp.Token

	return global.AdminToken, nil
}

func GetAllUsers(token string) ([]UserRecord, error) {
	url := os.Getenv("POCKETBASE_URL") + "/api/collections/users/records"
	page := 1

	var allItems []UserRecord

	for {
		query := map[string]string{
			"perPage": "500",
		}

		headers := map[string]string{
			"Authorization": token,
		}
		response, err := http.GetRequest(url, query, headers)
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
	url := os.Getenv("POCKETBASE_URL") + "/api/collections/songs/records"
	headers := map[string]string{
		"Authorization": token,
	}
	response, error := http.PostRequest(url, song, headers)
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
	url := os.Getenv("POCKETBASE_URL") + "/api/collections/users/records/" + user.ID
	header := map[string]string{
		"Authorization": token,
	}
	response, error := http.PatchRequest(url, user, header)
	if error != nil {
		return "", fmt.Errorf("failed to update user: %v", error)
	}
	return string(response), nil
}

func GetUserRecord(userID string, token string) (UserRecord, error) {
	if userID == "" {
		return UserRecord{}, fmt.Errorf("userID is empty")
	}
	url := os.Getenv("POCKETBASE_URL") + "/api/collections/users/records"
	headers := map[string]string{
		"Authorization": token,
	}

	query := map[string]string{
		"perPage": "1",
		"filter":  "(username='" + userID + "')",
	}
	response, err := http.GetRequest(url, query, headers)
	if err != nil {
		return UserRecord{}, fmt.Errorf("failed to get user record: %v", err)
	}

	userResponse := RecordsResponse{}
	err = json.Unmarshal(response, &userResponse)
	if err != nil {
		return UserRecord{}, fmt.Errorf("failed to unmarshal user response: %v", err)
	}

	if userResponse.TotalItems == 0 {
		return UserRecord{}, fmt.Errorf("user not found: %v", userID)
	}

	if userResponse.TotalItems > 1 {
		return UserRecord{}, fmt.Errorf("found more than one user with the same ID")
	}

	return userResponse.Items[0], nil
}

// CreatePlaylist creates a playlist in the database and returns the playlist ID.
func CreatePlaylist(playlist Playlist, token string) (string, error) {
	url := os.Getenv("POCKETBASE_URL") + "/api/collections/playlists/records"
	headers := map[string]string{
		"Authorization": token,
	}
	response, error := http.PostRequest(url, playlist, headers)
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
	url := os.Getenv("POCKETBASE_URL") + "/api/collections/song_playlist_link/records"
	headers := map[string]string{
		"Authorization": token,
	}
	body := SongPlaylistLink{
		PlaylistID: playlistID,
		SongID:     songID,
	}
	response, error := http.PostRequest(url, body, headers)
	if error != nil {
		return "", fmt.Errorf("failed to insert playlist song link: %v", error)
	}
	return string(response), nil
}

func GetSongBySongId(songID string, token string) (Song, error) {
	url := os.Getenv("POCKETBASE_URL") + "/api/collections/songs/records"
	headers := map[string]string{
		"Authorization": token,
	}

	query := map[string]string{
		"perPage": "1",
		"filter":  "(song_id='" + songID + "')",
	}
	response, err := http.GetRequest(url, query, headers)
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
	url := os.Getenv("POCKETBASE_URL") + "/api/collections/playlists/records"
	headers := map[string]string{
		"Authorization": token,
	}

	query := map[string]string{
		"perPage": "1",
		"filter":  "(hash='" + hash + "')",
	}
	response, err := http.GetRequest(url, query, headers)
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
	url := os.Getenv("POCKETBASE_URL") + "/api/collections/playlists/records"
	headers := map[string]string{
		"Authorization": token,
	}

	query := map[string]string{
		"perPage": "1000",
		"filter":  "(owner='" + userID + "')",
		"sort":    "-created",
	}
	response, err := http.GetRequest(url, query, headers)
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
	url := os.Getenv("POCKETBASE_URL") + "/api/collections/song_playlist_link/records"
	headers := map[string]string{
		"Authorization": token,
	}

	query := map[string]string{
		"perPage": "500",
		"filter":  "(playlist_id='" + playlistID + "')",
		"expand":  "song_id",
		"sort":    "created",
	}
	response, err := http.GetRequest(url, query, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist songs: %v", err)
	}
	// slog.Info("response", "response", string(response))

	songPlaylistResponse := PlaylistSongsResponse{}
	err = json.Unmarshal(response, &songPlaylistResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal song playlist response: %v", err)
	}

	var songs []Song
	for _, songPlaylistLink := range songPlaylistResponse.Items {
		songs = append(songs, songPlaylistLink.Expand.SongID)
	}

	return songs, nil
}

func ValidateToken(userId string, accessToken string, adminToken string) (bool, error) {
	user, err := GetUserRecord(userId, adminToken)
	if err != nil || user == (UserRecord{}) {
		slog.Error("Error getting user record: " + err.Error())
		return false, err
	}
	if user.AccessToken == accessToken {
		return true, nil
	} else {
		return false, nil
	}
}

func CreateUserRecord(record *spotify.PrivateUser, token *oauth2.Token, adminToken string) {
	url := os.Getenv("POCKETBASE_URL") + "/api/collections/users/records"
	headers := map[string]string{
		"Authorization": adminToken,
	}

	// literally does not matter what this is it just has to be unhackable 4head
	password := util.RandomString(25)

	// Create the request body
	data := map[string]string{
		"username":        record.ID,
		"email":           record.Email,
		"emailVisibility": "false",
		"verified":        "true",
		"password":        password,
		"passwordConfirm": password,
		"accessToken":     token.AccessToken,
		"refreshToken":    token.RefreshToken,
		"expiry":          util.FormatTime(token.Expiry),
		"display_name":    record.DisplayName,
		"avatar_url":      record.Images[0].URL,
	}

	_, error := http.PostRequest(url, data, headers)
	if error != nil {
		slog.Error("Error creating user record: " + error.Error())
	}
}

func GetPlaylistRecord(playlistId string, adminToken string) (Playlist, error) {
	url := os.Getenv("POCKETBASE_URL") + "/api/collections/playlists/records/" + playlistId
	// fmt.Println(url)
	headers := map[string]string{
		"Authorization": adminToken,
	}

	response, err := http.GetRequest(url, nil, headers)
	if err != nil {
		return Playlist{}, fmt.Errorf("failed to get playlist: %v", err)
	}

	playlistResponse := Playlist{}
	err = json.Unmarshal(response, &playlistResponse)
	if err != nil {
		return Playlist{}, fmt.Errorf("failed to unmarshal playlist response: %v", err)
	}

	return playlistResponse, nil
}
