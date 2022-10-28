package convert

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func extractHeader(msg *Message, content io.Reader) error {
	raw, err := io.ReadAll(content)
	if err != nil {
		return fmt.Errorf("reading content: %w", err)
	}

	lines := bytes.Split(raw, []byte("\n"))
	headerDividerCount := 0

	for _, line := range lines {
		if headerDividerCount == 2 {
			break
		}

		switch {
		case bytes.HasPrefix(line, []byte("---")):
			headerDividerCount++
		case bytes.HasPrefix(line, []byte("From:")):
			msg.From = string(bytes.TrimPrefix(line, []byte("From: ")))
		case bytes.HasPrefix(line, []byte("To:")):
			msg.To = string(bytes.TrimPrefix(line, []byte("To: ")))
		case bytes.HasPrefix(line, []byte("Subject:")):
			msg.Subject = string(bytes.TrimPrefix(line, []byte("Subject: ")))
		default:
			return fmt.Errorf("invalid header line: %s", line)
		}
	}

	return nil
}

func extractBody(msg *Message, content io.Reader) error {
	raw, err := io.ReadAll(content)
	if err != nil {
		return fmt.Errorf("reading content: %w", err)
	}

	lines := bytes.Split(raw, []byte("\n"))
	headerDividerCount := 0

	builder := strings.Builder{}

	for _, line := range lines {
		if headerDividerCount == 2 {
			builder.Write(line)
			builder.WriteString("\n")

			continue
		}

		if bytes.HasPrefix(line, []byte("---")) {
			headerDividerCount++
		}
	}

	switch {
	case headerDividerCount == 0:
		return fmt.Errorf("no header divider found")
	case headerDividerCount == 1:
		return fmt.Errorf("only one header divider found")
	case builder.Len() == 0:
		return errMissingBody
	}

	result := builder.String()
	result = strings.TrimPrefix(result, "\n")
	result = strings.TrimSpace(result)

	msg.Body = result

	return nil
}
