package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

const mimeTypeUserSessionsV1 = "application/vnd.userSessions.v1+json"

func (s *Service) getActiveSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Accept") {
	case "", mimeTypeUserSessionsV1: // including if not specified
		s.getActiveSessionsV1(w, r)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
}

func (s *Service) getActiveSessionsV1(w http.ResponseWriter, r *http.Request) {
	sessions, err := s.storage.GetSessionsV1(r.Context(), authenticatedUser(r))

	if err == nil {
		w.Header().Set("Content-Type", mimeTypeUserSessionsV1)
		err = json.NewEncoder(w).Encode(sessions)
	}

	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("Failed to get active sessions (v1, issue: %s).", requestId(r)))

		addIssueHeader(w, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
