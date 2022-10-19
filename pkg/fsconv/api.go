package fsconv

import (
	"bytes"
	"fmt"
	"path"

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
