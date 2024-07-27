package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"github.com/zmb3/spotify/v2"

	"daylist-rewind-backend/src/handlers"
	"daylist-rewind-backend/src/pocketbase"
	"daylist-rewind-backend/src/spotifyutil"
	"daylist-rewind-backend/src/util"
)

var adminToken string = ""

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	// Run scheduled tasks in the background
	go JobScheduler()

	// Serve the API on the main thread idk if this should be the other way around xd
	// mux := http.NewServeMux()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`This ain't for the best My reputation's never been worse, so You must like me for me We can't make Any promises now, can we, babe? But you can make me a drink`))
	})

	// Auth routes
	r.Get("/login", handlers.LoginHandler)
	r.Get("/callback", handlers.CallbackHandler)
	r.Get("/validate", handlers.ValidateToken)

	r.Route("/user/playlists", func(r chi.Router) {
		r.Get("/{userID}", handlers.GetUserPlaylistsHandler)
	})

	r.Route("/playlist", func(r chi.Router) {
		r.Get("/{playlistID}", handlers.GetSingleUserPlaylistHandler)
		r.Get("/{playlistID}/songs", handlers.GetPlaylistTracksHandler)
	})

	// Start the server on port 8080
	fmt.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func JobScheduler() {
	c := cron.New()
	// Poll all users for daylist updates every 30 seconds
	// Note: 30 seconds is spotify's rolling API rate limit window
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

	// generate a hash for the playlist to check if it already exists
	playlistHash := util.GetMD5Hash(
		// For hash collisions between users with the same playlist name
		user.ID +
			// Hash collisions between playlists with the same name
			daylist.Name +
			// Hash collisions between playlists that occur different days
			time.Now().Format("09-07-2017") +
			// Handles the edge case where the date changes at 12am but the playlist doesnt
			daylist.Tracks.Tracks[0].Track.ID.String(),
	)

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
