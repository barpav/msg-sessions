package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (s *Service) traceInternalServerError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err := errors.New(fmt.Sprintf("Recovered from panic: %v", rec))
				logAndReturnErrorWithIssue(w, r, err, "Internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

type authenticatedUserId struct{}

func (s *Service) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var authenticated bool
		var err error

		userId, password, parsed := r.BasicAuth()

		if parsed {
			authenticated, err = s.auth.ValidateCredentials(r.Context(), userId, password)
		}

		if err != nil {
			logAndReturnErrorWithIssue(w, r, err, "Credentials authentication failed")
			return
		}

		if !authenticated {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), authenticatedUserId{}, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func logAndReturnErrorWithIssue(w http.ResponseWriter, r *http.Request, err error, logMsg string) {
	issue := requestId(r)

	log.Err(err).Msg(fmt.Sprintf("Issue %s: %s", issue, logMsg))

	w.Header()["issue"] = []string{issue} // lowercase - non-canonical (vendor) header
	w.WriteHeader(http.StatusInternalServerError)
}

func authenticatedUser(r *http.Request) (id string) {
	id, _ = r.Context().Value(authenticatedUserId{}).(string)
	return id
}

func requestId(r *http.Request) string {
	id := r.Header.Get("request-id") // set by api-gateway

	if id != "" {
		return id
	}

	return "untraced"
}

func userIP(r *http.Request) string {
	ip := r.Header.Get("remote-addr") // set by api-gateway

	if ip != "" {
		return ip
	}

	return "unknown"
}

func userAgent(r *http.Request) string {
	agent := r.Header.Get("User-Agent")

	if agent != "" {
		return agent
	}

	return "unknown"
}
