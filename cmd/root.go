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
	log      *logrus.Logger
	fs       = &afero.Afero{Fs: afero.NewOsFs()}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "fssmtp",
	Short:        "A brief description of your application",
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var err error
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.fssmtp.yaml)")

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "log level [debug, info]")
	err = viper.BindPFlag(config.LogLevel, rootCmd.PersistentFlags().Lookup("log-level"))
	cobra.CheckErr(err)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".fssmtp" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(targetDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".fsmail")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	log = logging.GetLogger()
}
