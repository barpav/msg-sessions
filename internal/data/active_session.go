package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type ActiveSession struct {
	Key string
}

var ErrActiveSessionNotFound = errors.New("active session not found")

func (m *ActiveSession) Id(ctx context.Context, s *Storage) (id int64, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Failed to get active session id: %w", err)
		}
	}()

	if m.Key == "" {
		return 0, errors.New("session key must be specified")
	}

	id, err = s.db.HGet(ctx, sessionKeySum(m.Key), "id").Int64()

	switch {
	case err != nil:
		return 0, err
	case id == 0:
		return 0, ErrActiveSessionNotFound
	default:
		return id, nil
	}
}

func (m *ActiveSession) User(ctx context.Context, s *Storage) (id string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Failed to get active session user: %w", err)
		}
	}()

	if m.Key == "" {
		return "", errors.New("session key must be specified")
	}

	id, err = s.db.HGet(ctx, sessionKeySum(m.Key), "user").Result()

	switch {
	case err != nil:
		return "", err
	case id == "":
		return "", ErrActiveSessionNotFound
	default:
		return id, nil
	}
}

func (m *ActiveSession) Info(ctx context.Context, s *Storage) (info *SessionInfo, err error) {
	return info, err
}

func (m *ActiveSession) Delete(ctx context.Context, s *Storage) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Failed to delete active session: %w", err)
		}
	}()

	if m.Key == "" {
		return errors.New("session key must be specified")
	}

	var user string
	user, err = m.User(ctx, s)

	if err != nil {
		return err
	}

	var id int64
	id, err = m.Id(ctx, s)

	if err != nil {
		return err
	}

	var l *lock
	l, err = s.lock(ctx, sessionsLockKey(user), 30*time.Second)

	if err != nil {
		return err
	}

	defer func() {
		if err := l.unlock(); err != nil {
			log.Err(err).Msg("Failed to unlock user sessions.")
		}
	}()

	pipe := s.db.Pipeline()

	pipe.Del(ctx, sessionKeySum(m.Key), sessionInfoKey(user, id))

	pipe.SRem(ctx, sessionsIdsKey(user), id)

	pipe.Decr(ctx, sessionsTotalKey(user))

	commands, err := pipe.Exec(ctx)

	if err == nil {
		for _, cmd := range commands {
			err = errors.Join(err, cmd.Err())
		}
	}

	return err
}
