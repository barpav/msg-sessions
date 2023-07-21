package pb

import (
	"context"
	"fmt"
	"time"

	ssgrpc "github.com/barpav/msg-sessions/sessions_service_go_grpc"
	"github.com/rs/zerolog/log"
)

type ErrSessionNotFound interface {
	Error() string
	ImplementsSessionNotFoundError()
}

func (s *Service) Validate(ctx context.Context, sessionData *ssgrpc.SessionData) (*ssgrpc.ValidationResult, error) {
	userId, sessionId, err := s.storage.SessionKeyInfo(ctx, sessionData.Key)

	if err != nil {
		if _, ok := err.(ErrSessionNotFound); !ok {
			log.Err(err).Msg("Session data validation failed.")
			return nil, fmt.Errorf("Session data validation failed: %w", err)
		}
	}

	result := &ssgrpc.ValidationResult{User: userId}

	if userId != "" {
		go func() {
			info := map[string]interface{}{
				"lastActivity": time.Now().UTC(),
				"lastIp":       sessionData.Ip,
				"lastAgent":    sessionData.Agent,
			}

			err := s.storage.UpdateSessionInfo(context.Background(), userId, sessionId, info)

			if err != nil {
				log.Err(err).Msg(fmt.Sprintf("Failed to update user '%s' session '%d' info.", userId, sessionId))
			}
		}()
	}

	return result, nil
}
