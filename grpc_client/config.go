package sessions

import (
	"os"
	"strconv"
)

const (
	defaultHost        = "localhost"
	defaultPort        = "9000"
	defaultConnTimeout = 5
	defaultConnRetries = 5
)

const (
	envVarHost        = "MSG_SESSIONS_HOST"
	envVarPort        = "MSG_SESSIONS_PORT"
	envVarConnTimeout = "MSG_SESSIONS_CONN_TIMEOUT"
	envVarConnRetries = "MSG_SESSIONS_CONN_RETRIES"
)

type Config struct {
	host        string
	port        string
	connTimeout int
	connRetries int
}

func (c *Config) Read() {
	readSetting(envVarHost, defaultHost, &c.host)
	readSetting(envVarPort, defaultPort, &c.port)
	readNumericSetting(envVarConnTimeout, defaultConnTimeout, &c.connTimeout)
	readNumericSetting(envVarConnRetries, defaultConnRetries, &c.connRetries)
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
