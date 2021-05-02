package config

import (
	"os"

	"github.com/influenzanet/data-service/internal/constants"
)

// Config is the structure that holds all global configuration data
type Config struct {
	Port        string
	ServiceURLs struct {
		LoggingService string
		StudyService   string
	}
}

func InitConfig() Config {
	conf := Config{}
	conf.Port = os.Getenv(constants.ENV_DATA_SERVICE_LISTEN_PORT)
	conf.ServiceURLs.LoggingService = os.Getenv(constants.ENV_ADDR_LOGGING_SERVICE)
	conf.ServiceURLs.StudyService = os.Getenv(constants.ENV_ADDR_STUDY_SERVICE)
	return conf
}
