package email

import (
	"crypto/sha256"
	"fmt"
	"io"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"gopkg.in/gomail.v2"
)

func FetchInbox(log logger, credentials Credentials) ([]Message, error) {
	log.Debug("Connecting to IMAP server")

	client, err := client.DialTLS(credentials.IMAPServerAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("dialing: %w", err)
	}

	defer func() {
		_ = client.Logout()
	}()

	log.Debug("Logging in")

	if err = client.Login(credentials.Username, credentials.Password); err != nil {
		return nil, fmt.Errorf("logging in: %w", err)
	}

	inbox, err := client.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("selecting INBOX: %w", err)
	}

	if inbox.Messages == 0 {
		return nil, nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(inbox.Messages)

	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)

	log.Debug("Initiating fetch")

	go func() {
		if err := client.Fetch(seqset, items, messages); err != nil {
			fmt.Println("Fetch error:", err)
		}
	}()

	convertedMessages, err := handleMessages(section, messages)
	if err != nil {
		return nil, fmt.Errorf("handling messages: %w", err)
	}

	return convertedMessages, nil
}

func SendMessages(log logger, credentials Credentials, messages []Message) ([]string, error) {
	host, port, err := parseServerAddress(credentials.SMTPServerAddress)
	if err != nil {
		return nil, fmt.Errorf("parsing server address: %w", err)
	}

	dialer := gomail.NewDialer(host, port, credentials.Username, credentials.Password)

	sender, err := dialer.Dial()
	if err != nil {
		return nil, fmt.Errorf("dialing: %w", err)
	}

	receipts := make([]string, 0, len(messages))

	for _, message := range messages {
		m := gomail.NewMessage()
		m.SetHeader("From", message.From)
		m.SetHeader("To", message.To)
		m.SetHeader("Subject", message.Subject)

		rawBody, err := io.ReadAll(message.Body)
		if err != nil {
			return nil, fmt.Errorf("reading message body: %w", err)
		}

		m.SetBody("text/html", string(rawBody))

		if err := gomail.Send(sender, m); err != nil {
			return nil, fmt.Errorf("sending message: %w", err)
		}

		receipts = append(receipts, CalculateReceipt(message.From, message.To, message.Subject, string(rawBody)))

		time.Sleep(1 * time.Second)
	}

	return receipts, nil
}

func CalculateReceipt(from, to, subject, body string) string {
	hash := sha256.New()

	hash.Write([]byte(from))
	hash.Write([]byte(to))
	hash.Write([]byte(subject))
	hash.Write([]byte(body))

	return string(hash.Sum(nil))
}
