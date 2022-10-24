package cmd

import (
	"fmt"
	"os"

	"github.com/deifyed/fsmail/pkg/config"
	"github.com/deifyed/fsmail/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logLevel string
	cfgFile  string
	log      = &logrus.Logger{}
	fs       = &afero.Afero{Fs: afero.NewOsFs()}
)

var rootCmd = &cobra.Command{
	Use:          "fsmail",
	Short:        "fsmail enables you to synchronize your local directory with your email account",
	SilenceUsage: true,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var err error
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.fssmtp.yaml)")

	viper.SetDefault(config.LogLevel, "info")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", viper.GetString(config.LogLevel), "log level [debug, info]")
	err = viper.BindPFlag(config.LogLevel, rootCmd.PersistentFlags().Lookup("log-level"))
	cobra.CheckErr(err)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(targetDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".fsmail")
	}

	viper.AutomaticEnv()

	var msg string

	if err := viper.ReadInConfig(); err == nil {
		msg = fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed())
	} else {
		msg = "No config file found"
	}

	err := logging.ConfigureLogger(log, viper.GetString(config.LogLevel))
	cobra.CheckErr(err)

	log.Debug(msg)
}
