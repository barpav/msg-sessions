package pb

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/barpav/msg-sessions/internal/data"
	ssgrpc "github.com/barpav/msg-sessions/sessions_service_go_grpc"
	"google.golang.org/grpc"
)

type Service struct {
	Shutdown chan struct{}
	cfg      *Config
	server   *grpc.Server
	storage  Storage

	ssgrpc.UnimplementedSessionsServer
}

type Storage interface {
	SessionKeyInfo(ctx context.Context, key string) (userId string, sessionId int64, err error)
	UpdateSessionInfo(ctx context.Context, userId string, sessionId int64, info map[string]interface{}) (err error)
	EndAllSessions(ctx context.Context, userId string) (err error)
}

func (s *Service) Start(storage *data.Storage) {
	s.cfg = &Config{}
	s.cfg.Read()

	s.server = grpc.NewServer()
	s.storage = storage

	s.Shutdown = make(chan struct{}, 1)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.port))

		if err == nil {
			ssgrpc.RegisterSessionsServer(s.server, s)
			err = s.server.Serve(lis)
		}

		if err != nil {
			log.Err(err).Msg("gRPC server crashed.")
		}

		s.Shutdown <- struct{}{}
	}()
}

func (s *Service) Stop(ctx context.Context) (err error) {
	closed := make(chan struct{}, 1)

	go func() {
		s.server.GracefulStop()
		closed <- struct{}{}
	}()

	select {
	case <-closed:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("failed to stop gRPC service: %w", ctx.Err())
	}
}
