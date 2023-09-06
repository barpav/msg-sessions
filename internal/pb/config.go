package pb

import "os"

const defaultPort = "9000"

const envVarPort = "MSG_GRPC_PORT"

type config struct {
	port string
}

func (c *config) Read() {
	readSetting(envVarPort, defaultPort, &c.port)
}

func readSetting(setting, defaultValue string, result *string) {
	*result = os.Getenv(setting)
	if *result == "" {
		*result = defaultValue
	}
}
