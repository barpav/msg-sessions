package data

import (
	"os"
	"strconv"
)

const (
	defaultHost               = "localhost"
	defaultPort               = "6379"
	defaultMaxSessionsPerUser = 20
)

const (
	envVarHost               = "MSG_STORAGE_HOST"
	envVarPort               = "MSG_STORAGE_PORT"
	envVarMaxSessionsPerUser = "MSG_MAX_SESSIONS_PER_USER"
)

type Config struct {
	host               string
	port               string
	maxSessionsPerUser int
}

func (c *Config) Read() {
	readSetting(envVarHost, defaultHost, &c.host)
	readSetting(envVarPort, defaultPort, &c.port)
	readNumericSetting(envVarMaxSessionsPerUser, defaultMaxSessionsPerUser, &c.maxSessionsPerUser)
}

func readSetting(setting, defaultValue string, result *string) {
	*result = os.Getenv(setting)
	if *result == "" {
		*result = defaultValue
	}
}

func readNumericSetting(setting string, defaultValue int, result *int) {
	val := os.Getenv(setting)

	if val != "" {
		valNum, err := strconv.Atoi(val)

		if err == nil {
			*result = valNum
			return
		}
	}

	*result = defaultValue
}
