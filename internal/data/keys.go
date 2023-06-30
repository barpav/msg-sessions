package data

import "fmt"

func sessionsLockKey(user string) string {
	return fmt.Sprintf("user:%s:sessions:lock", user)
}

func sessionsTotalKey(user string) string {
	return fmt.Sprintf("user:%s:sessions:total", user)
}

func sessionsLastIdKey(user string) string {
	return fmt.Sprintf("user:%s:sessions:lastId", user)
}

func sessionInfoKey(user string, id int64) string {
	return fmt.Sprintf("user:%s:session:%d", user, id)
}

func sessionsIdsKey(user string) string {
	return fmt.Sprintf("user:%s:sessions", user)
}
