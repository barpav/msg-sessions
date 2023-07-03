package rest

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/barpav/msg-sessions/internal/rest/models"
)

type Service struct {
	Shutdown chan struct{}
	server   *http.Server
	auth     Authenticator
	storage  Storage
}

type Authenticator interface {
	ValidateCredentials(ctx context.Context, userId, password string) (valid bool, err error)
}

type Storage interface {
	StartNewSession(ctx context.Context, userId, ip, agent string) (id int64, key string, err error)
	GetSessionsV1(ctx context.Context, userId string) (sessions *models.UserSessionsV1, err error)
	EndSession(ctx context.Context, userId string, sessionId int64) (err error)
	EndAllSessions(ctx context.Context, userId string) (err error)

	SessionKeyInfo(ctx context.Context, key string) (userId string, sessionId int64, err error)
	UpdateSessionInfo(ctx context.Context, userId string, sessionId int64, info map[string]interface{}) (err error)
}

func (s *Service) Start(auth Authenticator, sessions Storage) {
	s.auth, s.storage = auth, sessions

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: s.operations(),
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

// Specification: https://barpav.github.io/msg-api-spec/#/sessions
func (s *Service) operations() *chi.Mux {
	ops := chi.NewRouter()

	ops.Use(s.traceInternalServerError)
	ops.Use(s.authenticate)

	ops.Post("/", s.startNewSession)
	ops.Get("/", s.getActiveSessions)
	ops.Delete("/", s.endSessions)

	return ops
}
