package spotifyutil

import (
	"context"
	"fmt"
	"log/slog"
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
