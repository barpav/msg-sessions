package data

import (
	"context"
	"fmt"
)

func (s *Storage) UpdateSessionInfo(ctx context.Context, userId string, sessionId int64, info map[string]interface{}) (err error) {
	err = s.db.HSet(ctx, sessionInfoKey(userId, sessionId), info).Err()

	if err != nil {
		return fmt.Errorf("Failed to update session info: %w", err)
	}

	return nil
}
