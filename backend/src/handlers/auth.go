package handlers

import (
	"crypto/sha256"
	"daylist-rewind-backend/src/pocketbase"
	"daylist-rewind-backend/src/util"
	"encoding/base64"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"sync"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

var (
	redirectURI       = "http://localhost:8080/callback"
	state             = "asdf"
	auth              = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadEmail, spotifyauth.ScopePlaylistModifyPrivate))
	codeVerifierStore = sync.Map{}
)

func generateCodeVerifier() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~"
	const length = 43

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateCodeChallenge(verifier string) string {
	h := sha256.New()
	h.Write([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadEmail, spotifyauth.ScopePlaylistModifyPrivate))
	codeVerifier := generateCodeVerifier()
	codeChallenge := generateCodeChallenge(codeVerifier)
	codeVerifierStore.Store(state, codeVerifier)

	authURL := auth.AuthURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)

	slog.Info("Redirecting to: " + authURL)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the code verifier using the state
	codeVerifierInterface, ok := codeVerifierStore.Load(state)
	if !ok {
		http.Error(w, "No code verifier found", http.StatusForbidden)
		slog.Error("No code verifier found")
	}
	codeVerifier, ok := codeVerifierInterface.(string)
	if !ok {
		http.Error(w, "Invalid code verifier format", http.StatusForbidden)
		slog.Error("Invalid code verifier format")
	}

	token, err := auth.Token(r.Context(), state, r, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		slog.Error(err.Error())
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		slog.Error("State mismatch: %s != %s\n", st, state)
	}

	client := spotify.New(auth.Client(r.Context(), token))
	user, err := client.CurrentUser(r.Context())
	if err != nil {
		http.Error(w, "Couldn't get user", http.StatusForbidden)
		slog.Error(err.Error())
	}

	// Display user info as a simple example
	w.Write([]byte("Login successful!"))
	w.Write([]byte("User ID: " + user.ID))
	w.Write([]byte("\nUser Name: " + user.DisplayName))

	admin_token, err := pocketbase.Authenticate(os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		log.Fatal("Error authenticating: " + err.Error())
	}

	record, err := pocketbase.GetUserRecord(user.ID, admin_token)
	if err != nil {
		slog.Error("Error getting user record: " + err.Error())
	}
	// user exists. update their auth token
	if record != (pocketbase.UserRecord{}) {
		token, err := client.Token()
		if err != nil {
			slog.Error("Error getting token: " + err.Error())
		}
		record.AccessToken = token.AccessToken
		record.RefreshToken = token.RefreshToken
		record.Expiry = util.FormatTime(token.Expiry)
		_, err = pocketbase.UpdateUser(record, admin_token)
		if err != nil {
			slog.Error("Error updating user record: " + err.Error())
		}
	} else {
		slog.Info("User does not exist")
		// user does not exist, create it
		// record = &pocketbase.UserRecord{
		// 	UserID: user.ID,
		// 	Email: user.Email,
		// 	AccessToken: token.AccessToken,
		// 	RefreshToken: token.RefreshToken,
		// 	Expiry: token.Expiry,
		// }
		// _, err = pocketbase.CreateUserRecord(*record, bearer)
		// if err != nil {
		// 	slog.Error("Error creating user record: " + err.Error())
		// }
	}

}
