package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/barpav/msg-sessions/internal/rest/models"
)

type Service struct {
	Shutdown chan struct{}
	cfg      *Config
	server   *http.Server
	auth     Authenticator
	storage  Storage
}

type Authenticator interface {
	ValidateCredentials(ctx context.Context, userId, password string) (valid bool, err error)
}

//go:generate mockery --name Storage
type Storage interface {
	StartNewSession(ctx context.Context, userId, ip, agent string) (id int64, key string, err error)
	GetSessionsV1(ctx context.Context, userId string) (sessions *models.UserSessionsV1, err error)
	EndSession(ctx context.Context, userId string, sessionId int64) (err error)
	EndAllSessions(ctx context.Context, userId string) (err error)
}

func (s *Service) Start(auth Authenticator, sessions Storage) {
	s.cfg = &Config{}
	s.cfg.Read()

	s.auth, s.storage = auth, sessions

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.cfg.port),
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
	err = s.server.Shutdown(ctx)

	if err != nil {
		err = fmt.Errorf("failed to stop HTTP service: %w", err)
	}

	return err
}

// Specification: https://barpav.github.io/msg-api-spec/#/sessions
func (s *Service) operations() *chi.Mux {
	ops := chi.NewRouter()

	ops.Use(s.traceInternalServerError)
	ops.Use(s.authenticate)

	// Public endpoint is the concern of the api gateway
	ops.Post("/", s.startNewSession)
	ops.Get("/", s.getActiveSessions)
	ops.Delete("/", s.endSessions)

	return ops
}
