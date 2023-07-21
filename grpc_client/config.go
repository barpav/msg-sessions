package sessions

import (
	"os"
	"strconv"
)

const (
	defaultHost        = "localhost"
	defaultPort        = "9000"
	defaultConnTimeout = 5
	defaultConnTries   = 5
)

const (
	envVarHost        = "MSG_SESSIONS_HOST"
	envVarPort        = "MSG_SESSIONS_PORT"
	envVarConnTimeout = "MSG_SESSIONS_CONN_TIMEOUT"
	envVarConnTries   = "MSG_SESSIONS_CONN_TRIES"
)

type Config struct {
	host        string
	port        string
	connTimeout int
	connTries   int
}

func (c *Config) Read() {
	readSetting(envVarHost, defaultHost, &c.host)
	readSetting(envVarPort, defaultPort, &c.port)
	readNumericSetting(envVarConnTimeout, defaultConnTimeout, &c.connTimeout)
	readNumericSetting(envVarConnTries, defaultConnTries, &c.connTries)

	if c.connTimeout < 1 {
		c.connTimeout = defaultConnTimeout
	}

	if c.connTries < 1 {
		c.connTries = 1
	}
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
