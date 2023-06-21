package api

import (
	"net/http"
)

type Service struct {
	// storage *data.Storage
}

func (s *Service) Init() error {
	// s.storage = &data.Storage{}
	// return s.storage.Init()
	return nil
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// URL handling is reverse proxy's concern
	switch r.Method {
	case http.MethodPost:
		s.startNewSession(w, r)
	case http.MethodGet:
		s.getActiveSessions(w, r)
	case http.MethodDelete:
		s.endSessions(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
