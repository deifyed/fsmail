package email

import (
	"fmt"
	"io"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
)

func handleMessages(section imap.BodySectionName, messages chan *imap.Message) ([]Message, error) {
	result := make([]Message, 0)

	for {
		msg, ok := <-messages
		if !ok {
			break
		}

		extractedMessage, err := extractMessage(&section, msg)
		if err != nil {
			return nil, fmt.Errorf("parsing message: %w", err)
		}

		result = append(result, Message{
			From:    extractedMessage.From,
			To:      extractedMessage.To,
			Subject: extractedMessage.Subject,
			Body:    extractedMessage.Body,
		})
	}

	return result, nil
}

func extractMessage(section *imap.BodySectionName, rawMessage *imap.Message) (Message, error) {
	resultMessage := Message{}

	r := rawMessage.GetBody(section)
	if r == nil {
		fmt.Println("Server didn't returned message body")
		resultMessage.Body = strings.NewReader("<!-- no content -->")
	}

	mailReader, err := mail.CreateReader(r)
	if err != nil {
		return Message{}, fmt.Errorf("creating mail reader: %w", err)
	}

	header := mailReader.Header

	if from, err := header.AddressList("From"); err == nil {
		resultMessage.From = from[0].Address
	}
	if to, err := header.AddressList("To"); err == nil {
		resultMessage.To = to[0].Address
	}
	if subject, err := header.Subject(); err == nil {
		resultMessage.Subject = subject
	}

	for {
		p, err := mailReader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return Message{}, fmt.Errorf("reading mail part: %w", err)
		}

		switch p.Header.(type) {
		case *mail.InlineHeader:
			resultMessage.Body = p.Body
		case *mail.AttachmentHeader:
			fmt.Println("Attachment found")
		}
	}

	return resultMessage, nil
}
