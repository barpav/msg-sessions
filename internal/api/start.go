package api

import (
	"fmt"
	"net/http"
)

func (s *Service) startNewSession(w http.ResponseWriter, r *http.Request) {
	user, password, ok := r.BasicAuth()

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Printf("user: %s, password: %s", user, password)
}
