package pb

import (
	"context"
	"fmt"

	ssgrpc "github.com/barpav/msg-sessions/sessions_service_go_grpc"
	"github.com/rs/zerolog/log"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) EndAll(ctx context.Context, user *ssgrpc.User) (*emptypb.Empty, error) {
	userId := user.GetId()
	err := s.storage.EndAllSessions(ctx, userId)

	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("Failed to end all user '%s' sessions.", userId))
		return nil, fmt.Errorf("failed to end all user sessions: %w", err)
	}

	log.Info().Msg(fmt.Sprintf("All user '%s' sessions ended.", userId))

	return &emptypb.Empty{}, nil
}
