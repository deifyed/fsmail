package sync

import (
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/deifyed/fssmtp/pkg/credentials"
	"github.com/deifyed/fssmtp/pkg/fsconv"
	"github.com/deifyed/fssmtp/pkg/keyring"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/gomail.v2"
)

func RunE(fs *afero.Afero, targetDir *string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		workDir, err := filepath.Abs(*targetDir)
		if err != nil {
			return fmt.Errorf("acquiring absolute target dir: %w", err)
		}

		creds, err := acquireCredentials()
		if err != nil {
			return fmt.Errorf("acquiring credentials: %w", err)
		}

		host, port, err := parseServerAddress(creds.ServerAddress)
		if err != nil {
			return fmt.Errorf("parsing server address: %w", err)
		}

		dialer := gomail.NewDialer(host, port, creds.Username, creds.Password)

		sender, err := dialer.Dial()
		if err != nil {
			return fmt.Errorf("dialing: %w", err)
		}

		defer func() {
			_ = sender.Close()
		}()

		err = handleOutbox(fs, sender, path.Join(workDir, "outbox"))
		if err != nil {
			return fmt.Errorf("handling outbox: %w", err)
		}

		return nil
	}
}

func handleOutbox(fs *afero.Afero, sender gomail.Sender, outboxDir string) error {
	messages, err := fsconv.DirectoryToMessages(fs, outboxDir)
	if err != nil {
		return fmt.Errorf("extracting messages: %w", err)
	}

	preparedMessages := make([]*gomail.Message, len(messages))

	for index, message := range messages {
		m := gomail.NewMessage()

		m.SetHeader("From", "fssmtp@localhost")
		m.SetHeader("To", message.Recipient)
		m.SetHeader("Subject", message.Subject)
		m.SetBody("text/html", message.Body)

		preparedMessages[index] = m
	}

	err = gomail.Send(sender, preparedMessages...)
	if err != nil {
		return fmt.Errorf("sending: %w", err)
	}

	return nil
}

func acquireCredentials() (credentials.Credentials, error) {
	store := keyring.Client{Prefix: generatePrefix("")}

	creds := credentials.Credentials{}
	var err error

	creds.ServerAddress, err = store.Get(credentials.CredentialsSecretName, credentials.ServerAddressKey)
	if err != nil {
		return credentials.Credentials{}, fmt.Errorf("retrieving server address: %w", err)
	}

	creds.Username, err = store.Get(credentials.CredentialsSecretName, credentials.UsernameKey)
	if err != nil {
		return credentials.Credentials{}, fmt.Errorf("retrieving username: %w", err)
	}

	creds.Password, err = store.Get(credentials.CredentialsSecretName, credentials.PasswordKey)
	if err != nil {
		return credentials.Credentials{}, fmt.Errorf("retrieving password: %w", err)
	}

	return creds, nil
}

func parseServerAddress(serverAddress string) (string, int, error) {
	parts := strings.Split(serverAddress, ":")

	host := parts[0]

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("converting port from string to int: %w", err)
	}

	return host, port, nil
}

func generatePrefix(username string) string {
	return "fssmtp"
}
