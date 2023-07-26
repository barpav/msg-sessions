package rest

import (
	"fmt"
	"net/http"
)

type ErrTooManySessions interface {
	Error() string
	ImplementsTooManySessionsError()
}

// https://barpav.github.io/msg-api-spec/#/sessions/post_sessions
func (s *Service) startNewSession(w http.ResponseWriter, r *http.Request) {
	id, key, err := s.storage.StartNewSession(r.Context(), authenticatedUser(r), userIP(r), userAgent(r))

	if err != nil {
		if _, ok := err.(ErrTooManySessions); ok {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		logAndReturnErrorWithIssue(w, r, err, "Failed to start new session")
		return
	}

	// lowercase - non-canonical (vendor) headers
	w.Header()["session-id"] = []string{fmt.Sprintf("%d", id)}
	w.Header()["session-key"] = []string{key}

	w.WriteHeader(http.StatusCreated)
}
