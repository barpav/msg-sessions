package rest

import (
	"encoding/json"
	"net/http"
)

const mimeTypeUserSessionsV1 = "application/vnd.userSessions.v1+json"

// https://barpav.github.io/msg-api-spec/#/sessions/get_sessions
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
	sessions, err := s.storage.UserSessionsV1(r.Context(), authenticatedUser(r))

	if err == nil {
		w.Header().Set("Content-Type", mimeTypeUserSessionsV1)
		err = json.NewEncoder(w).Encode(sessions)
	}

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, "Failed to get active sessions (v1)")
		return
	}

	w.WriteHeader(http.StatusOK)
}
