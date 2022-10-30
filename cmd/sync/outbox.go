package sync

import (
	"fmt"
	"os"
	"path"
	"strings"

	stdfs "io/fs"

	"github.com/deifyed/fsmail/pkg/convert"
	"github.com/deifyed/fsmail/pkg/credentials"
	"github.com/deifyed/fsmail/pkg/email"
	"github.com/deifyed/fsmail/pkg/keyring"
	"github.com/spf13/afero"
)

func handleOutbox(log logger, fs *afero.Afero, absoluteOutboxDirectory string, absoluteSentDirectory string, creds credentials.Credentials) error {
	files, err := fs.ReadDir(absoluteOutboxDirectory)
	if err != nil {
		return fmt.Errorf("reading outbox directory: %w", err)
	}

	files = filterFiles(files)

	messages := make([]email.Message, len(files))
	receiptMap := make(map[string]string)

	for index, file := range files {
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

		receiptMap[email.CalculateReceipt(msg.From, msg.To, msg.Subject, msg.Body)] = filename
		messages[index] = convertMessageToEmail(msg)
	}

	receipts, err := email.SendMessages(log, email.Credentials{
		SMTPServerAddress: creds.SMTPServerAddress,
		Username:          creds.Username,
		Password:          creds.Password,
	}, messages)
	if err != nil {
		log.Warn(fmt.Errorf("sending messages: %w", err).Error())
	}

	if len(receipts) == 0 {
		return err
	}

	for _, receipt := range receipts {
		filename := receiptMap[receipt]
		sourcePath := path.Join(absoluteOutboxDirectory, filename)
		destinationPath := path.Join(absoluteSentDirectory, filename)

		err = fs.Rename(sourcePath, destinationPath)
		if err != nil {
			return fmt.Errorf("moving file: %w", err)
		}
	}

	return err
}

func convertMessageToEmail(msg convert.Message) email.Message {
	return email.Message{
		From:    msg.From,
		To:      msg.To,
		Subject: msg.Subject,
		Body:    strings.NewReader(msg.Body),
	}
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

func generatePrefix(username string) string {
	return "fssmtp"
}
