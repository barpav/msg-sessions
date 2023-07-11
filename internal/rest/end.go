package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

func (s *Service) endSessions(w http.ResponseWriter, r *http.Request) {
	var err error
	if id := r.URL.Query().Get("id"); id != "" {
		sessionId, _ := strconv.ParseInt(id, 0, 64)
		err = s.storage.EndSession(r.Context(), authenticatedUser(r), sessionId)
	} else {
		err = s.storage.EndAllSessions(r.Context(), authenticatedUser(r))
	}

	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("Failed to end sessions (issue: %s).", requestId(r)))

		addIssueHeader(w, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
