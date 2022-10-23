package logging

import (
	"errors"
	"fmt"
	"os"

	"github.com/deifyed/fsmail/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var errInvalidLevel = errors.New("invalid level")

func ConfigureLogger(log *logrus.Logger) error {
	var err error

	log.Out = os.Stdout
	log.Formatter = &logrus.JSONFormatter{PrettyPrint: true}

	log.Level, err = parseLevel(viper.GetString(config.LogLevel))
	if err != nil {
		return fmt.Errorf("parsing log level: %w", err)
	}

	return nil
}

func parseLevel(level string) (logrus.Level, error) {
	switch level {
	case "debug":
		return logrus.DebugLevel, nil
	case "info":
		return logrus.InfoLevel, nil
	case "warn":
		return logrus.WarnLevel, nil
	case "error":
		return logrus.ErrorLevel, nil
	default:
		return logrus.InfoLevel, errInvalidLevel
	}
}
