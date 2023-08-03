package lyricsify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"net/http"

	"github.com/ahmagdy/lyricsify/config"
	"github.com/rs/cors"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/zap"
)

const (
	_authServerPort = ":8080"
	// redirectURI is the OAuth redirect URI for the application.
	// You must register an application at Spotify's developer portal
	// and enter this value.
	_redirectURI = "http://localhost%s/callback"
	_state       = "abc123"

	_ctxTimeout = 5 * time.Second
)

type authServer struct {
	logger *zap.Logger

	httpServer    *http.Server
	auth          *spotifyauth.Authenticator
	authCompleted chan struct{}
	spotifyClient *spotify.Client
}

func NewAuthServer(logger *zap.Logger, cfg *config.Config) *authServer {
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(fmt.Sprintf(_redirectURI, _authServerPort)),
		spotifyauth.WithScopes(spotifyauth.ScopeUserLibraryRead, spotifyauth.ScopeUserReadPrivate),
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

func (a *authServer) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/callback", a.completeAuthHandler)
	mux.HandleFunc("/sync", a.handleSync)
	mux.HandleFunc("/search", a.handleSearch)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Got request for", zap.String("reqURL", r.URL.String()))
	})

	handler := cors.AllowAll().Handler(mux)

	a.httpServer = &http.Server{Addr: _authServerPort, Handler: handler}
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("error return while shutting down the server", zap.Error(err))
		}
	}()

	a.logger.Info("startnig the server on ", zap.String("serverAddr", fmt.Sprintf("http://localhost%s/sync", _authServerPort)))

	return nil
}

func (a *authServer) getAuthURL() string {
	return a.auth.AuthURL(_state)
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
	token, err := a.auth.Token(r.Context(), _state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		a.logger.Fatal("Couldn't get token", zap.Error(err))
	}

	if st := r.FormValue("state"); st != _state {
		http.NotFound(w, r)
		a.logger.Fatal("Couldn't get token", zap.String("receivedState", st), zap.String("originalState", _state))
	}

	// use the token to get an authenticated client
	client := spotify.New(a.auth.Client(r.Context(), token))
	// fmt.Fprintf(w, "Login Completed!")

	a.spotifyClient = client
	a.authCompleted <- struct{}{}

	http.Redirect(w, r, "http://localhost:3000/", http.StatusSeeOther)
}

type msg struct {
	AuthURL string `json:"auth_url"`
}

func (a *authServer) handleSync(w http.ResponseWriter, r *http.Request) {
	resp := &msg{AuthURL: a.getAuthURL()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		a.logger.Error("failed to encode json response", zap.Error(err))
	}
}

func (a *authServer) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	svc, err := New(r.Context(), a.SpotifyClient())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		a.logger.Error("failed to init svc", zap.Error(err))
		return
	}

	resp, err := svc.Search(r.Context(), query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		a.logger.Error("failed to init svc", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		a.logger.Error("failed to encode json response", zap.Error(err))
	}
}

// WaitForAuthToBeCompleted ..
func (a *authServer) WaitForAuthToBeCompleted() {
	<-a.authCompleted
}

// Stop graceful shutdown
// You can use somethinglike
// // Setting up signal capturing
// stop := make(chan os.Signal, 1)
// signal.Notify(stop, os.Interrupt)

// // Waiting for SIGINT (kill -2)
// <-stop
func (a *authServer) Stop(ctx context.Context) error {
	if a.httpServer == nil {
		return nil
	}

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("httpServer.Shutdown: %w", err)
	}

	return nil
}
