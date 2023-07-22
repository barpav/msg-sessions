package models

// Schema: userSessions.v1
type UserSessionsV1 struct {
	Active int              `json:"active"`
	List   []*UserSessionV1 `json:"list"`
}

type UserSessionV1 struct {
	Id           int64   `json:"id"`
	Created      UtcTime `json:"created"`
	LastActivity UtcTime `json:"lastActivity"`
	LastIp       string  `json:"lastIp"`
	LastAgent    string  `json:"lastAgent"`
}
