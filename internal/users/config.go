package users

import (
	"os"
	"strconv"
)

const (
	defaultHost        = "localhost"
	defaultPort        = "9000"
	defaultConnTimeout = 30
)

const (
	envVarHost        = "MSG_USERS_HOST"
	envVarPort        = "MSG_USERS_PORT"
	envVarConnTimeout = "MSG_USERS_CONN_TIMEOUT"
)

type Config struct {
	host        string
	port        string
	connTimeout int
}

func (c *Config) Read() {
	readSetting(envVarHost, defaultHost, &c.host)
	readSetting(envVarPort, defaultPort, &c.port)
	readNumericSetting(envVarConnTimeout, defaultConnTimeout, &c.connTimeout)
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
