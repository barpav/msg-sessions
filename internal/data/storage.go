package data

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Storage struct {
	cfg *Config
	db  *redis.Client
}

func (s *Storage) Open() error {
	s.cfg = &Config{}
	s.cfg.Read()
	return s.connectToDatabase()
}

func (s *Storage) Close(ctx context.Context) (err error) {
	var closeErr error
	closed := make(chan struct{}, 1)

	go func() {
		closeErr = s.db.Close()
		closed <- struct{}{}
	}()

	select {
	case <-closed:
		err = closeErr
	case <-ctx.Done():
		err = ctx.Err()
	}

	if err != nil {
		err = fmt.Errorf("failed to disconnect from database: %w", err)
	}

	return err
}

func (s *Storage) connectToDatabase() error {
	s.db = redis.NewClient(&redis.Options{
		Addr: net.JoinHostPort(s.cfg.host, s.cfg.port),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.db.Ping(ctx).Err()

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Info().Msg(fmt.Sprintf("Successfully connected to database at %s", s.db.Options().Addr))

	return nil
}
