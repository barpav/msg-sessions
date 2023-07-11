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
	closed := make(chan struct{}, 1)

	go func() {
		err = s.db.Close()

		closed <- struct{}{}
	}()

	select {
	case <-closed:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Storage) connectToDatabase() error {
	s.db = redis.NewClient(&redis.Options{
		Addr: net.JoinHostPort(s.cfg.host, s.cfg.port),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.db.Ping(ctx).Err()

	if err != nil {
		return fmt.Errorf("Failed to connect to database: %s", err)
	}

	log.Info().Msg(fmt.Sprintf("Successfully connected to database at %s", s.db.Options().Addr))

	return nil
}
