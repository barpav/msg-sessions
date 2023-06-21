package data

import "os"

const (
	defaultHost     = "localhost"
	defaultPort     = "5432"
	defaultDatabase = "postgres"
	defaultUser     = "postgres"
	defaultPassword = "postgres"
)

const (
	envVarHost     = "MSG_STORAGE_HOST"
	envVarPort     = "MSG_STORAGE_PORT"
	envVarDatabase = "MSG_STORAGE_DATABASE"
	envVarUser     = "MSG_STORAGE_USER"
	envVarPassword = "MSG_STORAGE_PASSWORD"
)

type Config struct {
	host     string
	port     string
	database string
	user     string
	password string
}

func (c *Config) Read() {
	readSetting(envVarHost, defaultHost, &c.host)
	readSetting(envVarPort, defaultPort, &c.port)
	readSetting(envVarDatabase, defaultDatabase, &c.database)
	readSetting(envVarUser, defaultUser, &c.user)
	readSetting(envVarPassword, defaultPassword, &c.password)
}

func readSetting(setting, defaultValue string, result *string) {
	*result = os.Getenv(setting)
	if *result == "" {
		*result = defaultValue
	}
}
