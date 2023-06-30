package data

import (
	"crypto/md5"
	"fmt"
)

func sessionKeySum(key string) string {
	sum := md5.Sum([]byte(key))
	return fmt.Sprintf("%x", sum)
}
