package rest

import "net/http"

func (s *Service) getActiveSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Accept") {
	case "", "application/vnd.userSessions.v1+json": // by default
		s.getActiveSessionsV1(w, r)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
}

func (s *Service) getActiveSessionsV1(w http.ResponseWriter, r *http.Request) {

}
