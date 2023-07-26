package rest

import (
	"net/http"
	"strconv"
)

// https://barpav.github.io/msg-api-spec/#/sessions/delete_sessions
func (s *Service) endSessions(w http.ResponseWriter, r *http.Request) {
	var err error
	if id := r.URL.Query().Get("id"); id != "" {
		sessionId, _ := strconv.ParseInt(id, 0, 64)
		err = s.storage.EndSession(r.Context(), authenticatedUser(r), sessionId)
	} else {
		err = s.storage.EndAllSessions(r.Context(), authenticatedUser(r))
	}

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, "Failed to end sessions")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
