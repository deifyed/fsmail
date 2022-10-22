package cmd

import (
	"fmt"
	"os"

	"github.com/deifyed/fsmail/cmd/sync"
	"github.com/deifyed/fsmail/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	imapServerAddress string
	smtpServerAddress string
	targetDir         string
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "synchronizes directory with server",
	RunE:  sync.RunE(log, fs, &targetDir),
}

func init() {
	rootCmd.AddCommand(syncCmd)

	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("getting work directory: %w", err))
	}

	syncCmd.Flags().StringVarP(&targetDir, "directory", "d", workDir, "target directory")
	err = viper.BindPFlag(config.WorkingDirectory, syncCmd.Flags().Lookup("directory"))
	cobra.CheckErr(err)

	syncCmd.Flags().StringVarP(&imapServerAddress, "imap-server-address", "i", "", "IMAP server address")
	err = viper.BindPFlag(config.IMAPServerAddress, syncCmd.Flags().Lookup("imap-server-address"))
	cobra.CheckErr(err)

	syncCmd.Flags().StringVarP(&smtpServerAddress, "smtp-server-address", "s", "", "SMTP server address")
	err = viper.BindPFlag(config.SMTPServerAddress, syncCmd.Flags().Lookup("smtp-server-address"))
	cobra.CheckErr(err)
}
