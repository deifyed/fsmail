package fsconv

import (
	"bytes"
	_ "embed"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/spf13/afero"
)

func DirectoryToMessages(fs *afero.Afero, targetDir string) ([]Message, error) {
	files, err := fs.ReadDir(targetDir)
	if err != nil {
		return nil, fmt.Errorf("listing: %w", err)
	}

	messages := make([]Message, 0)

	for _, file := range files {
		filePath := path.Join(targetDir, file.Name())

		raw, err := fs.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("reading: %w", err)
		}

		delim := []byte("\n")
		parts := bytes.Split(raw, delim)

		subject := parts[0]
		body := bytes.Join(parts[1:], delim)

		messages = append(messages, Message{Recipient: file.Name(), Subject: string(subject), Body: string(body)})
	}

	return messages, nil
}

func MessagesToDirectory(fs *afero.Afero, targetDir string, messages []Message) error {
	err := fs.MkdirAll(targetDir, 0o755)
	if err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	return nil
}

//go:embed file-template.md
var messageFileTemplate string

func messageToFile(fs *afero.Afero, targetDir string, message Message) error {
	t, err := template.New("message").Parse(messageFileTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	buf := bytes.Buffer{}

	err = t.Execute(&buf, struct {
		To      string
		Cc      []string
		Bcc     []string
		Subject string
		Body    string
	}{
		To:      message.Recipient,
		Cc:      message.Cc,
		Subject: message.Subject,
		Body:    message.Body,
	})
	if err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	err = fs.WriteReader(path.Join(targetDir, subjectAsFilename(message.Subject)), &buf)
	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

func subjectAsFilename(subject string) string {
	return strings.ReplaceAll(subject, " ", "-")
}
