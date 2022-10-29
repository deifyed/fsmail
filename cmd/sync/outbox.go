package sync

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	stdfs "io/fs"

	"github.com/deifyed/fsmail/pkg/convert"
	"github.com/deifyed/fsmail/pkg/credentials"
	"github.com/deifyed/fsmail/pkg/keyring"
	"github.com/spf13/afero"
	"gopkg.in/gomail.v2"
)

func handleOutbox(fs *afero.Afero, absoluteOutboxDirectory string, absoluteSentDirectory string, creds credentials.Credentials) error {
	sender, err := getSender(creds)
	if err != nil {
		return fmt.Errorf("getting sender: %w", err)
	}

	defer func() {
		_ = sender.Close()
	}()

	files, err := fs.ReadDir(absoluteOutboxDirectory)
	if err != nil {
		return fmt.Errorf("reading outbox directory: %w", err)
	}

	files = filterFiles(files)

	for _, file := range files {
		filename := file.Name()

		filePath := path.Join(absoluteOutboxDirectory, filename)

		f, err := fs.OpenFile(filePath, os.O_RDONLY, defaultFilePermissions)
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}

		msg, err := convert.ToMessage(f)
		if err != nil {
			return fmt.Errorf("converting file to message: %w", err)
		}

		err = gomail.Send(sender, convertMessageToGoMail(msg))
		if err != nil {
			return fmt.Errorf("sending message: %w", err)
		}

		err = fs.Rename(filePath, path.Join(absoluteSentDirectory, filename))
		if err != nil {
			return fmt.Errorf("moving file: %w", err)
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

func convertMessageToGoMail(msg convert.Message) *gomail.Message {
	m := gomail.NewMessage()

	m.SetHeader("From", msg.From)
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/plain", msg.Body)

	return m
}

const defaultFilePermissions = 0o600

func filterFiles(files []stdfs.FileInfo) []stdfs.FileInfo {
	var filteredFiles []stdfs.FileInfo

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filteredFiles = append(filteredFiles, file)
	}

	return filteredFiles
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
