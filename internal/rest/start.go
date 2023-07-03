package rest

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ErrTooManySessions interface {
	Error() string
	ImplementsTooManySessionsError()
}

func (s *Service) startNewSession(w http.ResponseWriter, r *http.Request) {
	id, key, err := s.storage.StartNewSession(r.Context(), authenticatedUser(r), userIP(r), userAgent(r))

	if err != nil {
		if _, ok := err.(ErrTooManySessions); ok {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		log.Err(err).Msg(fmt.Sprintf("Failed to start new session (issue: %s).", requestId(r)))

		addIssueHeader(w, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// lowercase - non-canonical (vendor) headers
	w.Header()["session-id"] = []string{fmt.Sprintf("%d", id)}
	w.Header()["session-key"] = []string{key}

	w.WriteHeader(http.StatusCreated)
}
