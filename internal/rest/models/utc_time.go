package models

import (
	"fmt"
	"time"
)

type UtcTime time.Time

func (t UtcTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(t).UTC().Format("2006-01-02T15:04:05Z"))
	return []byte(formatted), nil
}
