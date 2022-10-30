package sync

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/deifyed/fsmail/pkg/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunE(log logger, fs *afero.Afero, targetDir *string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		absoluteWorkDirectory, err := filepath.Abs(*targetDir)
		if err != nil {
			return fmt.Errorf("acquiring absolute target dir: %w", err)
		}

		absoluteInboxDirectory := path.Join(absoluteWorkDirectory, "inbox")
		absoluteOutboxDirectory := path.Join(absoluteWorkDirectory, "outbox")
		absoluteSentDirectory := path.Join(absoluteWorkDirectory, "sent")

		imapServerAddress := viper.GetString(config.IMAPServerAddress)
		smtpServerAddress := viper.GetString(config.SMTPServerAddress)

		log.Debugf("Using work dir: %s", absoluteWorkDirectory)
		log.Debugf("Using IMAP server address: %s", imapServerAddress)
		log.Debugf("Using SMTP server address: %s", smtpServerAddress)

		log.Debug("Preparing credentials")

		creds, err := acquireCredentials(imapServerAddress, smtpServerAddress)
		if err != nil {
			return fmt.Errorf("acquiring credentials: %w", err)
		}

		err = handleInbox(log, fs, absoluteInboxDirectory, creds)
		if err != nil {
			return fmt.Errorf("handling inbox: %w", err)
		}

		err = handleOutbox(log, fs, absoluteOutboxDirectory, absoluteSentDirectory, creds)
		if err != nil {
			return fmt.Errorf("handling outbox: %w", err)
		}

		return nil
	}
}
