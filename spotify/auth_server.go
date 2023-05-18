package spotify

import (
	"context"
	"fmt"

	"net/http"

	"github.com/ahmagdy/lyricsify/config"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/zap"
)

const (
	// redirectURI is the OAuth redirect URI for the application.
	// You must register an application at Spotify's developer portal
	// and enter this value.
	redirectURI = "http://localhost:8080/callback"
	state       = "abc123"
)

type authServer struct {
	logger *zap.Logger

	auth          *spotifyauth.Authenticator
	authCompleted chan struct{}
	spotifyClient *spotify.Client
}

func NewAuthServer(logger *zap.Logger, cfg *config.Config) *authServer {
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
		spotifyauth.WithClientID(cfg.SpotifyID),
		spotifyauth.WithClientSecret(cfg.SpotifySecret),
	)
	return &authServer{
		logger:        logger,
		auth:          auth,
		authCompleted: make(chan struct{}),
	}
}

func (a *authServer) SpotifyClient() *spotify.Client {
	return a.spotifyClient
}

func (a *authServer) Start() {
	// first start an HTTP server
	http.HandleFunc("/callback", a.completeAuthHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Got request for", zap.String("reqURL", r.URL.String()))
	})

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			a.logger.Fatal("failed to start the http server", zap.Error(err))
		}
	}()

	url := a.auth.AuthURL(state)
	a.logger.Info("Please log in to Spotify by visiting the following page in your browser", zap.String("reqURL", url))
}

// use the client to make calls that require authorization
func (a *authServer) Verify(ctx context.Context) error {
	user, err := a.spotifyClient.CurrentUser(ctx)
	if err != nil {
		return fmt.Errorf("spotifyClient.CurrentUser: %w", err)
	}

	a.logger.Info("You are logged in as", zap.String("userID", user.ID))
	return nil
}

func (a *authServer) completeAuthHandler(w http.ResponseWriter, r *http.Request) {
	tok, err := a.auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		a.logger.Fatal("Couldn't get token", zap.Error(err))
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		a.logger.Fatal("Couldn't get token", zap.String("receivedState", st), zap.String("originalState", state))
	}

	// use the token to get an authenticated client
	client := spotify.New(a.auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")

	a.spotifyClient = client
	a.authCompleted <- struct{}{}
}

// WaitForAuthToBeCompleted ..
func (a *authServer) WaitForAuthToBeCompleted() {
	<-a.authCompleted
}
