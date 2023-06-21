package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/barpav/msg-sessions/internal/data"
	"github.com/barpav/msg-sessions/internal/users"
)

type Service struct {
	Shutdown chan struct{}
	server   *http.Server
	storage  *data.Storage
	users    *users.Client
}

func (s *Service) Start(storage *data.Storage, users *users.Client) {
	s.storage, s.users = storage, users

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: s,
	}

	s.Shutdown = make(chan struct{}, 1)

	go func() {
		err := s.server.ListenAndServe()

		if err != http.ErrServerClosed {
			log.Err(err).Msg("HTTP server crashed.")
		}

		s.Shutdown <- struct{}{}
	}()
}

func (s *Service) Stop(ctx context.Context) (err error) {
	return s.server.Shutdown(ctx)
}

// https://barpav.github.io/msg-api-spec/#/sessions
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authenticated, err := s.authenticated(r)

	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("Authentication failed (issue: %s).", r.Header.Get("request-id")))

		w.Header()["issue"] = []string{r.Header.Get("request-id")} // lowercase - non-canonical (vendor) header
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !authenticated {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// URL handling is reverse proxy's concern
	switch r.Method {
	case http.MethodPost:
		s.startNewSession(w, r)
	case http.MethodGet:
		s.getActiveSessions(w, r)
	case http.MethodDelete:
		s.endSessions(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Service) authenticated(r *http.Request) (ok bool, err error) {
	userId, password, parsed := r.BasicAuth()

	if parsed {
		ok, err = s.users.ValidateCredentials(r.Context(), userId, password)
	}

	return ok, err
}
