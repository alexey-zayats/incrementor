package config

import (
	"fmt"
	"github.com/gravitational/trace"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config struct holds whole application configuration
type Config struct {
	SQL struct {
		Driver   string
		Database string
		Username string
		Password string
		Hostname string
		Port     string
		SslMode  string
	}
	Log struct {
		Level string
	}
	Debug  bool
	Server struct {
		Listen struct {
			Network string
			Address string
		}
		TLS struct {
			CertFile string
			KeyFile  string
		}
	}
	Client struct {
		Auth struct {
			Username string
			Password string
		}
		TLS struct {
			CertFile string
		}
		Dial struct {
			Address string
		}
	}
	JWT struct {
		Secret   string
		Duration string
	}
	Incrementor struct {
		MaxValue    int
		MinValue    int
		IncrementBy int
	}
}

// NewConfig read config from file or ENV and unmarshal to to Config struct
func NewConfig() (*Config, error) {

	var level logrus.Level
	var err error
	c := &Config{}

	if err = viper.ReadInConfig(); err == nil {
		logrus.Infof("Using config file: %s", viper.ConfigFileUsed())
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, errors.Wrap(err, "Error with viper.Unmarshal")
	}

	if c.Debug {
		logrus.SetReportCaller(true)
	}

	fmt.Printf("%v\n", c)

	level, err = logrus.ParseLevel(c.Log.Level)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error with logrus.ParseLevel(%s)", c.Log.Level))
	}

	logrus.SetLevel(level)

	if udphook, err := trace.NewUDPHook(); err == nil {
		logrus.AddHook(udphook)
	}

	return c, nil
}
