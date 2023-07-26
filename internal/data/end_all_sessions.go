package data

import (
	"context"
	"fmt"
	"strconv"
)

func (s *Storage) EndAllSessions(ctx context.Context, userId string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to end all sessions: %w", err)
		}
	}()

	var ids []string
	ids, err = s.db.SMembers(ctx, sessionsIdsKey(userId)).Result()

	switch {
	case err != nil:
		return err
	case len(ids) == 0: // operation is idempotent
		return nil
	}

	var sessionId int64
	for _, id := range ids {
		sessionId, err = strconv.ParseInt(id, 0, 64)

		if err == nil {
			err = s.EndSession(ctx, userId, sessionId)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
