package models

import "time"

type UserSessionsV1 struct {
	Active int
	List   []*UserSessionV1
}

type UserSessionV1 struct {
	Id           int64
	Created      time.Time
	LastActivity time.Time
	LastIp       string
	LastAgent    string
}
