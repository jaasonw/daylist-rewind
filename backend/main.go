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
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
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

	bearer, err := pocketbase.Authenticate(os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		log.Fatal("Error authenticating" + err.Error())
	}

	items, err := pocketbase.GetAllUsers(bearer)
	if err != nil {
		log.Fatal("Error getting users" + err.Error())
	}
	for _, userRecord := range items {

		client, error := spotifyutil.CreateClient(userRecord.AccessToken, userRecord.RefreshToken, userRecord.Expiry)
		if error != nil {
			slog.Error("Error creating client" + error.Error())
			continue
		}
		token, err := client.Token()
		if err != nil {
			// increment error count
			userRecord.ErrorCount++
			userRecord.LastError = err.Error()
			pocketbase.UpdateUser(userRecord, bearer)
			slog.Error("Error getting token" + err.Error())
		}

		userRecord.AccessToken = token.AccessToken
		userRecord.RefreshToken = token.RefreshToken
		userRecord.Expiry = util.FormatTime(token.Expiry)
		_, err = pocketbase.UpdateUser(userRecord, bearer)
		if err != nil {
			slog.Error("Error updating user" + err.Error())
			continue
		}

		user, err := client.CurrentUser(context.Background())
		if err != nil {
			slog.Error("Error getting user" + err.Error())
			continue
		}
		daylist, err := client.GetPlaylist(context.Background(), spotify.ID("37i9dQZF1EP6YuccBxUcC1"))
		if err != nil {
			slog.Error("Error getting playlist" + err.Error())
			continue
		}
		slog.Info("Processing daylist for:", "user.ID", user.ID)
		slog.Info("Playlist:", "daylist.Name", daylist.Name)
		daylistJson, _ := json.Marshal(daylist)
		slog.Info("Playlist:", "daylist.hash", util.GetMD5Hash(string(daylistJson)))
		playlistExists, _, _ := pocketbase.CheckPlaylistExists(util.GetMD5Hash(string(daylistJson)), bearer)
		if playlistExists {
			slog.Info("Daylist unchanged, skipping")
			continue
		}
		playlistId, err := pocketbase.CreatePlaylist(pocketbase.Playlist{
			Hash:        util.GetMD5Hash(string(daylistJson)),
			Title:       daylist.Name,
			Owner:       userRecord.ID,
			Description: daylist.Description,
		}, bearer)

		slog.Info("Created Playlist:", "playlist.ID", playlistId)

		if err != nil {
			slog.Error("Error creating playlist" + err.Error())
			continue
		}
		for _, track := range daylist.Tracks.Tracks {
			// slog.Info(
			// 	"Processing",
			// 	"track",
			// 	track.Track.Name+" - "+
			// 		track.Track.Artists[0].Name,
			// )
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
			insertResponse, err := pocketbase.InsertSong(song, bearer)
			if err != nil {
				if insertResponse != "Value must be unique." {
					log.Printf("Error inserting song" + song.SongID + err.Error())
				} else {
					log.Printf("Song already exists: " + song.SongID)
					// find it
					song, err := pocketbase.GetSongBySongId(song.SongID, bearer)
					if err != nil {
						slog.Error("Error getting song by song id" + err.Error())
					}
					insertResponse = song.ID
				}
			}
			_, err = pocketbase.AddSongToPlaylist(playlistId, insertResponse, bearer)
			if err != nil {
				slog.Error("Error adding song to playlist" + err.Error())
			}

		}
	}
}
