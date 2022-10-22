package logging

import (
	"errors"
	"os"

	"github.com/deifyed/fsmail/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var errInvalidLevel = errors.New("invalid level")

func GetLogger() *logrus.Logger {
	var err error

	log := logrus.New()

	log.Out = os.Stdout
	log.Formatter = &logrus.JSONFormatter{PrettyPrint: true}

	log.Level, err = parseLevel(viper.GetString(config.LogLevel))
	if err != nil {
		if !errors.Is(err, errInvalidLevel) {
			panic(err.Error())
		}

		log.Warnf("defaulting to %s upon invalid log level \"%s\"", log.Level, viper.GetString(config.LogLevel))
	}

	return log
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
