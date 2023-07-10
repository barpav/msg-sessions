package data

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/barpav/msg-sessions/internal/rest/models"
	"github.com/redis/go-redis/v9"
)

func (s *Storage) GetSessionsV1(ctx context.Context, userId string) (sessions *models.UserSessionsV1, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Failed to get sessions (v1): %w", err)
		}
	}()

	var ids []string
	ids, err = s.db.SMembers(ctx, sessionsIdsKey(userId)).Result()

	if err != nil {
		return nil, err
	}

	sessions = &models.UserSessionsV1{
		Active: len(ids),
		List:   make([]*models.UserSessionV1, 0, len(ids)),
	}

	var (
		pipe                                                 redis.Pipeliner
		createdCmd, lastActivityCmd, lastIpCmd, lastAgentCmd *redis.StringCmd
		t                                                    time.Time
	)
	for _, id := range ids {
		session := &models.UserSessionV1{}

		session.Id, err = strconv.ParseInt(id, 0, 64)

		if err != nil {
			return nil, err
		}

		pipe = s.db.Pipeline()

		createdCmd = s.db.HGet(ctx, sessionInfoKey(userId, session.Id), "created")
		lastActivityCmd = s.db.HGet(ctx, sessionInfoKey(userId, session.Id), "lastActivity")
		lastIpCmd = s.db.HGet(ctx, sessionInfoKey(userId, session.Id), "lastIp")
		lastAgentCmd = s.db.HGet(ctx, sessionInfoKey(userId, session.Id), "lastAgent")

		_, err = pipe.Exec(ctx)

		if err != nil {
			return nil, err
		}

		t, err = createdCmd.Time()

		if err != nil {
			return nil, err
		}

		session.Created = models.UtcTime(t)

		t, err = lastActivityCmd.Time()

		if err != nil {
			return nil, err
		}

		session.LastActivity = models.UtcTime(t)

		session.LastIp, err = lastIpCmd.Result()

		if err != nil {
			return nil, err
		}

		session.LastAgent, err = lastAgentCmd.Result()

		if err != nil {
			return nil, err
		}

		sessions.List = append(sessions.List, session)
	}

	return sessions, nil
}
