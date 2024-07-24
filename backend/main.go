package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"github.com/zmb3/spotify/v2"

	"daylist-rewind-backend/pocketbase"
	"daylist-rewind-backend/spotifyutil"
	"daylist-rewind-backend/util"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Run scheduled tasks in the background
	go JobScheduler()

	mux := http.NewServeMux()

	// Register the helloHandler function for the "/" route
	mux.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(responseWriter, "Hello, World!")
	})

	mux.Handle("/playlist/", http.StripPrefix("/playlist/", http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		// Extract the playlist ID from the URL path
		path := strings.Trim(request.URL.Path, "/")
		playlistID := path // The path is already stripped by http.StripPrefix

		if playlistID == "" {
			http.Error(responseWriter, "Missing playlist ID", http.StatusBadRequest)
			return
		}

		bearer, err := pocketbase.Authenticate(os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD"))
		if err != nil {
			log.Fatal("Error authenticating" + err.Error())
		}
		songs, err := pocketbase.GetPlaylistSongs(playlistID, bearer)
		if err != nil {
			slog.Error("Error getting playlist songs" + err.Error())
			fmt.Fprintf(responseWriter, "Error getting playlist songs: %s", err)
		}

		// Process the playlist ID (for example, lookup playlist details)
		fmt.Fprintf(responseWriter, "Playlist ID: %s", playlistID)
		for _, song := range songs {
			songJson, _ := json.Marshal(song)
			fmt.Fprintf(responseWriter, "Song: %s\n", songJson)
		}
	})))

	// Start the server on port 8080
	fmt.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func JobScheduler() {
	c := cron.New()
	// run updateuser every 30 seconds
	c.AddFunc("*/30 * * * * *", func() {
		bearer, err := pocketbase.Authenticate(os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD"))
		if err != nil {
			log.Fatal("Error authenticating" + err.Error())
		}
		UpdateUsers(bearer)
	})

	// Jobs to add later
	// 1. clean up old playlists
	// 2. garbage collect unreferenced songs
	c.Start()
}

func UpdateUsers(pocketbase_token string) {
	// Retrive list of users from the database
	items, err := pocketbase.GetAllUsers(pocketbase_token)
	if err != nil {
		log.Fatal("Error getting users" + err.Error())
	}
	// Shuffle it to kinda spread out the api calls lol idk
	util.ShuffleArray(items)

	for _, userRecord := range items {
		// Authenticate with Spotify
		client, error := spotifyutil.CreateClient(userRecord.AccessToken, userRecord.RefreshToken, userRecord.Expiry)
		if error != nil {
			slog.Error("Error creating client" + error.Error())
			continue
		}
		token, err := client.Token()
		if err != nil {
			slog.Error("Error getting token" + err.Error())
		}

		// Update the auth token in the database
		userRecord.AccessToken = token.AccessToken
		userRecord.RefreshToken = token.RefreshToken
		userRecord.Expiry = util.FormatTime(token.Expiry)
		_, err = pocketbase.UpdateUser(userRecord, pocketbase_token)
		if err != nil {
			slog.Error("Error updating user" + err.Error())
			continue
		}
		UpdateUser(client, userRecord, pocketbase_token)

		// Sleep for n seconds to avoid rate limiting
		time.Sleep(10 * time.Second)
	}
}

func UpdateUser(client *spotify.Client, userRecord pocketbase.UserRecord, pocketbase_token string) error {
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		slog.Error("Error getting user" + err.Error())
		return err
	}
	daylist, err := client.GetPlaylist(context.Background(), spotify.ID("37i9dQZF1EP6YuccBxUcC1"))
	if err != nil {
		slog.Error("Error getting playlist" + err.Error())
		return err
	}

	playlistHash := util.GetMD5Hash(daylist.Name + time.Now().Format("09-07-2017"))

	slog.Info("Processing daylist for:", "user.ID", user.ID)
	slog.Info("Playlist:", "daylist.Name", daylist.Name)
	slog.Info("Playlist:", "daylist.Image", daylist.Images[0].URL)
	slog.Info("Playlist:", "daylist.hash", playlistHash)
	playlistExists, _, _ := pocketbase.CheckPlaylistExists(playlistHash, pocketbase_token)
	if playlistExists {
		slog.Info("Daylist unchanged, skipping")
		return nil
	}
	playlistId, err := pocketbase.CreatePlaylist(pocketbase.Playlist{
		Hash:        playlistHash,
		Title:       daylist.Name,
		Owner:       userRecord.ID,
		Description: daylist.Description,
	}, pocketbase_token)

	slog.Info("Created Playlist:", "playlist.ID", playlistId)

	if err != nil {
		slog.Error("Error creating playlist" + err.Error())
		return err
	}
	for _, track := range daylist.Tracks.Tracks {
		song := pocketbase.Song{
			SongID:     track.Track.ID.String(),
			Name:       track.Track.Name,
			Artist:     track.Track.Artists[0].Name,
			ArtistURL:  track.Track.Artists[0].ExternalURLs["spotify"],
			Album:      track.Track.Album.Name,
			AlbumURL:   track.Track.Album.ExternalURLs["spotify"],
			AlbumCover: track.Track.Album.Images[0].URL,
			PreviewURL: track.Track.PreviewURL,
			Duration:   int(track.Track.Duration),
		}
		insertResponse, err := pocketbase.InsertSong(song, pocketbase_token)
		if err != nil {
			if insertResponse != "Value must be unique." {
				log.Printf("Error inserting song" + song.SongID + err.Error())
			} else {
				song, err := pocketbase.GetSongBySongId(song.SongID, pocketbase_token)
				if err != nil {
					slog.Error("Error getting song by song id" + err.Error())
				}
				insertResponse = song.ID
			}
		}
		_, err = pocketbase.AddSongToPlaylist(playlistId, insertResponse, pocketbase_token)
		if err != nil {
			slog.Error("Error adding song to playlist" + err.Error())
		}
	}
	return nil
}
