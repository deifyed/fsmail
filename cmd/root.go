package cmd

import (
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

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "log level [debug, info]")
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

	err := logging.ConfigureLogger(log)
	cobra.CheckErr(err)

	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}
}
