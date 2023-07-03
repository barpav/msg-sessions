package data

import (
	"context"

	"github.com/barpav/msg-sessions/internal/rest/models"
)

func (s *Storage) GetSessionsV1(ctx context.Context, userId string) (sessions *models.UserSessionsV1, err error) {
	return sessions, err
}
