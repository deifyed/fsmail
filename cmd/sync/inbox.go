package sync

import (
	"fmt"

	"github.com/deifyed/fsmail/pkg/credentials"
	"github.com/deifyed/fsmail/pkg/email"
	"github.com/deifyed/fsmail/pkg/fsconv"
	"github.com/spf13/afero"
)

func handleInbox(log logger, fs *afero.Afero, absoluteInboxDirectory string, creds credentials.Credentials) error {
	log.Debug("Fetching inbox messages")

	messages, err := email.FetchInbox(log, email.Credentials{
		IMAPServerAddress: creds.IMAPServerAddress,
		Username:          creds.Username,
		Password:          creds.Password,
	})
	if err != nil {
		return fmt.Errorf("fetching inbox: %w", err)
	}

	log.Debugf("Saving %d inbox messages to %s", len(messages), absoluteInboxDirectory)

	for _, msg := range messages {
		err = fsconv.WriteMessageToDirectory(fs, absoluteInboxDirectory, emailMessageToFsConvMessage(msg))
		if err != nil {
			return fmt.Errorf("writing message to directory: %w", err)
		}
	}

	return nil
}

func emailMessageToFsConvMessage(source email.Message) fsconv.Message {
	return fsconv.Message{
		From:    source.From,
		To:      source.To,
		Subject: source.Subject,
		Body:    source.Body,
	}
}
