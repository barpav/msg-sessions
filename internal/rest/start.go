package rest

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/barpav/msg-sessions/internal/data"
	"github.com/rs/zerolog/log"
)

func (s *Service) startNewSession(w http.ResponseWriter, r *http.Request) {
	user, ctx := authenticatedUser(r), r.Context()

	newSession := data.NewSession{User: user}
	createdSession, err := newSession.Create(ctx, s.storage)

	if errors.Is(err, data.ErrTooManySessions) {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("Failed to start new session (issue: %s).", requestId(r)))

		addIssueHeader(w, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	creationTime := time.Now()

	info := data.SessionInfo{
		User:         user,
		Id:           newSession.Id(),
		Created:      creationTime,
		LastActivity: creationTime,
		LastIp:       userIP(r),
		LastAgent:    userAgent(r),
	}

	err = info.Update(ctx, s.storage)

	if err != nil {
		reqId := requestId(r)

		log.Err(err).Msg(fmt.Sprintf("Failed to start new session: failed to update new session info (issue: %s).", reqId))

		err = createdSession.Delete(ctx, s.storage)

		if err == nil {
			log.Info().Msg(fmt.Sprintf("Inconsistent session successfully deleted (issue: %s).", reqId))
		} else {
			log.Err(err).Msg(fmt.Sprintf("Failed to delete inconsistent session (issue: %s).", reqId))
		}

		addIssueHeader(w, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// lowercase - non-canonical (vendor) headers
	w.Header()["session-id"] = []string{fmt.Sprintf("%d", newSession.Id())}
	w.Header()["session-key"] = []string{createdSession.Key}

	w.WriteHeader(http.StatusCreated)
}
