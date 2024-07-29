package spotifyutil

import (
	"context"
	"daylist-rewind-backend/src/pocketbase"
	"daylist-rewind-backend/src/util"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func CreateClient(access_token string, refresh_token string, expiry string) (*spotify.Client, error) {
	var auth = spotifyauth.New()

	layout := "2006-01-02 15:04:05.000Z"
	expiryTime, err := time.Parse(layout, expiry)
	if err != nil {
		slog.Error("Error parsing time" + err.Error())
		return nil, err
	}

	token := &oauth2.Token{
		AccessToken:  access_token,
		TokenType:    "Bearer",
		RefreshToken: refresh_token,
		Expiry:       expiryTime,
	}
	token, err = auth.RefreshToken(context.Background(), token)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to refresh token: %v", err))
		return nil, err
	}
	client := spotify.New(auth.Client(context.Background(), token))
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to get current user: %v", err))
		return nil, err
	}

	// Print the user info
	slog.Debug("Sucessfully authenticated: ", "user.ID", user.ID, "user.DisplayName", user.DisplayName)
	return client, nil
}

// Takes a Playlist ID (database id) and exports it to the user's spotify account
func ExportPlaylist(client *spotify.Client, playlistId string) (string, error) {

	bearer, err := pocketbase.Authenticate(os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		log.Fatal("Error authenticating: " + err.Error())
	}

	playlistRecord, err := pocketbase.GetPlaylistRecord(playlistId, bearer)
	if err != nil {
		slog.Error("Error getting playlist record: " + err.Error())
		return "", err
	}

	songs, err := pocketbase.GetPlaylistSongs(playlistId, bearer)
	if err != nil {
		slog.Error("Error getting playlist songs: " + err.Error())
		return "", err
	}
	// Create the playlist
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to get current user: %v", err))
		return "", err
	}

	playlist, err := client.CreatePlaylistForUser(
		context.Background(),
		user.ID,
		playlistRecord.Title,
		util.SanitizeHTML(playlistRecord.Description),
		false,
		false,
	)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create playlist: %v",
			err))
		return "", err
	}

	// Add the songs to the playlist
	var songIds []spotify.ID
	for _, song := range songs {
		songIds = append(songIds, spotify.ID(song.SongID))
	}

	_, err = client.AddTracksToPlaylist(context.Background(), playlist.ID, songIds...)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to add tracks to playlist: %v", err))
		return "", err
	}
	return playlist.ID.String(), nil
}
