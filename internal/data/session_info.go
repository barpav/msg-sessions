package data

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type SessionInfo struct {
	User         string
	Id           int64
	Created      time.Time
	LastActivity time.Time
	LastIp       string
	LastAgent    string
}

func (m *SessionInfo) Update(ctx context.Context, s *Storage) (err error) {
	if m.User == "" || m.Id == 0 {
		return errors.New("Failed to update session info: user and session id must be specified.")
	}

	fields := make(map[string]interface{}, 4)

	if !m.Created.IsZero() {
		fields["created"] = m.Created
	}

	if !m.LastActivity.IsZero() {
		fields["lastActivity"] = m.LastActivity
	}

	if m.LastIp != "" {
		fields["lastIp"] = m.LastIp
	}

	if m.LastAgent != "" {
		fields["lastAgent"] = m.LastAgent
	}

	err = s.db.HSet(ctx, sessionInfoKey(m.User, m.Id), fields).Err()

	if err != nil {
		return fmt.Errorf("Failed to update session info: %w", err)
	}

	return nil
}

func (m *SessionInfo) Total(ctx context.Context, s *Storage) (err error) {
	// s.db.SCard()
	return err
}

func (m *SessionInfo) Next(ctx context.Context, s *Storage) (err error) {
	// if Id == 0 -> from start

	return err
}
