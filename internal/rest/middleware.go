package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (s *Service) handleInternalServerError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err := errors.New(fmt.Sprintf("Recovered from panic: %v", rec))
				log.Err(err).Msg(fmt.Sprintf("Internal server error (issue: %s).", r.Header.Get("request-id")))

				w.Header()["issue"] = []string{r.Header.Get("request-id")} // lowercase - non-canonical (vendor) header
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (s *Service) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var authenticated bool
		var err error

		userId, password, parsed := r.BasicAuth()

		if parsed {
			authenticated, err = s.users.ValidateCredentials(r.Context(), userId, password)
		}

		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("Authentication failed (issue: %s).", r.Header.Get("request-id")))

			w.Header()["issue"] = []string{r.Header.Get("request-id")} // lowercase - non-canonical (vendor) header
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !authenticated {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
