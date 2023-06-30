package rest

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
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
