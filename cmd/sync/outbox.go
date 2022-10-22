package sync

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/deifyed/fsmail/pkg/credentials"
	"github.com/deifyed/fsmail/pkg/fsconv"
	"github.com/deifyed/fsmail/pkg/keyring"
	"github.com/spf13/afero"
	"gopkg.in/gomail.v2"
)

func handleOutbox(fs *afero.Afero, outboxDir string, creds credentials.Credentials) error {
	host, port, err := parseServerAddress(creds.SMTPServerAddress)
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

func acquireCredentials(imapServerAddress string, smtpServerAddress string) (credentials.Credentials, error) {
	var err error
	store := keyring.Client{Prefix: generatePrefix("")}

	creds := credentials.Credentials{
		IMAPServerAddress: imapServerAddress,
		SMTPServerAddress: smtpServerAddress,
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
