package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

var (
	// ConfigFile holds path to configuration file
	ConfigFile string
	// EnvPrefix holds prefix for ENV variables
	EnvPrefix = "INC"
)

// Init initialize configurations after cobra had inited
func Init() {
	logrus.SetLevel(logrus.TraceLevel)

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stderr)

	viper.SetConfigFile(ConfigFile)
}
