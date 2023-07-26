package data

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type lock struct {
	storage  *Storage
	ctx      context.Context
	resource string
	seal     string
}

func (s *Storage) lock(ctx context.Context, resource string, timeout time.Duration) (l *lock, err error) {
	l = &lock{
		storage:  s,
		ctx:      ctx,
		resource: resource,
		seal:     uuid.NewString(),
	}

	set, err := s.db.SetNX(ctx, l.resource, l.seal, timeout).Result()

	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to lock resource '%s': %w", l.resource, err)
	case !set:
		return nil, fmt.Errorf("failed to lock resource '%s': already locked.", l.resource)
	default:
		return l, nil
	}
}

func (l *lock) unlock() error {
	result, err := unlocking.Run(l.ctx, l.storage.db, []string{l.resource}, l.seal).Int()

	switch {
	case err != nil:
		return fmt.Errorf("failed to unlock resource '%s': %w", l.resource, err)
	case result == 0:
		return fmt.Errorf("failed to unlock resource '%s': lock not found.", l.resource)
	default:
		return nil
	}
}

// https://redis.io/docs/manual/patterns/distributed-locks/
var unlocking = redis.NewScript(`
if redis.call("get",KEYS[1]) == ARGV[1] then
	return redis.call("del",KEYS[1])
else
	return 0
end
`)
