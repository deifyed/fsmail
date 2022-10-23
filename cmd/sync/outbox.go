package sync

import (
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"

	"github.com/deifyed/fsmail/pkg/credentials"
	"github.com/deifyed/fsmail/pkg/fsconv"
	"github.com/deifyed/fsmail/pkg/keyring"
	"github.com/spf13/afero"
	"gopkg.in/gomail.v2"
)

func handleOutbox(fs *afero.Afero, outboxDir string, sentDir string, creds credentials.Credentials) error {
	sender, err := getSender(creds)
	if err != nil {
		return fmt.Errorf("getting sender: %w", err)
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

		rawBody, err := io.ReadAll(message.Body)
		if err != nil {
			return fmt.Errorf("buffering body: %w", err)
		}

		m.SetHeader("From", "fssmtp@localhost")
		m.SetHeader("To", message.To)
		m.SetHeader("Subject", message.Subject)
		m.SetBody("text/html", string(rawBody))

		preparedMessages[index] = m
	}

	err = gomail.Send(sender, preparedMessages...)
	if err != nil {
		return fmt.Errorf("sending: %w", err)
	}

	err = moveFiles(fs, outboxDir, sentDir)
	if err != nil {
		return fmt.Errorf("moving files: %w", err)
	}

	return nil
}

func getSender(creds credentials.Credentials) (gomail.SendCloser, error) {
	host, port, err := parseServerAddress(creds.SMTPServerAddress)
	if err != nil {
		return nil, fmt.Errorf("parsing server address: %w", err)
	}

	dialer := gomail.NewDialer(host, port, creds.Username, creds.Password)

	sender, err := dialer.Dial()
	if err != nil {
		return nil, fmt.Errorf("dialing: %w", err)
	}

	return sender, nil
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

func moveFiles(fs *afero.Afero, sourceDir string, destinationDir string) error {
	files, err := fs.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("listing outbox directory: %w", err)
	}

	for _, file := range files {
		src := path.Join(sourceDir, file.Name())
		dest := path.Join(destinationDir, file.Name())

		err = fs.Rename(src, dest)
		if err != nil {
			return fmt.Errorf("removing file: %w", err)
		}
	}

	return nil
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
