package sessions

import (
	"os"
)

const (
	defaultHost = "localhost"
	defaultPort = "9000"
)

const (
	envVarHost = "MSG_SESSIONS_HOST"
	envVarPort = "MSG_SESSIONS_PORT"
)

type Config struct {
	host string
	port string
}

func (c *Config) Read() {
	readSetting(envVarHost, defaultHost, &c.host)
	readSetting(envVarPort, defaultPort, &c.port)
}

func readSetting(setting, defaultValue string, result *string) {
	*result = os.Getenv(setting)
	if *result == "" {
		*result = defaultValue
	}
}
