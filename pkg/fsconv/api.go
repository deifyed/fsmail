package fsconv

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"
	"text/template"

	"github.com/spf13/afero"
)

type header struct {
	From    string
	To      string
	Subject string
}

func extractHeader(content io.Reader) (header, error) {
	raw, err := io.ReadAll(content)
	if err != nil {
		return header{}, fmt.Errorf("reading content: %w", err)
	}

	lines := bytes.Split(raw, []byte("\n"))
	headerDividerCount := 0

	hdr := header{}

	for _, line := range lines {
		if headerDividerCount == 2 {
			break
		}

		switch {
		case bytes.HasPrefix(line, []byte("---")):
			headerDividerCount++
		case bytes.HasPrefix(line, []byte("From:")):
			hdr.From = string(bytes.TrimPrefix(line, []byte("From: ")))
		case bytes.HasPrefix(line, []byte("To:")):
			hdr.To = string(bytes.TrimPrefix(line, []byte("To: ")))
		case bytes.HasPrefix(line, []byte("Subject:")):
			hdr.Subject = string(bytes.TrimPrefix(line, []byte("Subject: ")))
		default:
			return header{}, fmt.Errorf("invalid header line: %s", line)
		}
	}

	return hdr, nil
}

func extractBody(content io.Reader) (io.Reader, error) {
	raw, err := io.ReadAll(content)
	if err != nil {
		return nil, fmt.Errorf("reading content: %w", err)
	}

	lines := bytes.Split(raw, []byte("\n"))
	headerDividerCount := 0

	body := bytes.Buffer{}

	for _, line := range lines {
		if headerDividerCount == 2 {
			body.Write(line)
			body.Write([]byte("\n"))
			continue
		}

		switch {
		case bytes.HasPrefix(line, []byte("---")):
			headerDividerCount++
		default:
			continue
		}
	}

	bodyAsBytes := bytes.TrimSpace(body.Bytes())

	return bytes.NewBuffer(bodyAsBytes), nil
}

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

		hdr, err := extractHeader(bytes.NewReader(raw))
		if err != nil {
			return nil, fmt.Errorf("extracting header: %w", err)
		}

		body, err := extractBody(bytes.NewReader(raw))
		if err != nil {
			return nil, fmt.Errorf("extracting body: %w", err)
		}

		messages = append(messages, Message{
			To:      hdr.To,
			From:    hdr.From,
			Subject: hdr.Subject,
			Body:    body,
		})
	}

	return messages, nil
}

func WriteMessagesToDirectory(fs *afero.Afero, targetDir string, messages []Message) error {
	err := fs.MkdirAll(targetDir, 0o755)
	if err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	for _, message := range messages {
		err = WriteMessageToDirectory(fs, targetDir, message)
		if err != nil {
			return fmt.Errorf("writing message: %w", err)
		}
	}

	return nil
}

func WriteMessageToDirectory(fs *afero.Afero, targetDir string, message Message) error {
	t, err := template.New("message").Parse(messageFileTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	buf := bytes.Buffer{}

	rawBody, err := io.ReadAll(message.Body)
	if err != nil {
		return fmt.Errorf("buffering body: %w", err)
	}

	err = t.Execute(&buf, struct {
		To      string
		Cc      string
		Bcc     string
		Subject string
		Body    string
	}{
		To:      message.To,
		Cc:      formatList(message.Cc),
		Subject: message.Subject,
		Body:    string(rawBody),
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

func formatList(list []string) string {
	return strings.Join(list, ", ")
}
