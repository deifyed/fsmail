package convert

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

func ToMessage(content io.Reader) (Message, error) {
	contentCopy := bytes.Buffer{}
	tee := io.TeeReader(content, &contentCopy)

	msg := Message{}

	if err := extractHeader(&msg, tee); err != nil {
		return Message{}, fmt.Errorf("extracting header: %w", err)
	}

	if err := extractBody(&msg, &contentCopy); err != nil {
		if !errors.Is(err, errMissingBody) {
			return Message{}, fmt.Errorf("extracting body: %w", err)
		}

		msg.Body = "<!-- empty -->"
	}

	return msg, nil
}

func ToReader(msg Message) io.Reader {
	buf := bytes.Buffer{}

	buf.Write([]byte(divider + "\n"))

	buf.Write([]byte("To: " + msg.To + "\n"))
	buf.Write([]byte("From: " + msg.From + "\n"))
	buf.Write([]byte("Subject: " + msg.Subject + "\n"))

	buf.Write([]byte(divider + "\n\n"))

	buf.Write([]byte(msg.Body + "\n"))

	return &buf
}
