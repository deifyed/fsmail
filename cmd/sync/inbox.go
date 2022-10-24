package sync

import (
	"fmt"
	"io"
	"strings"

	"github.com/deifyed/fsmail/pkg/credentials"
	"github.com/deifyed/fsmail/pkg/fsconv"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/spf13/afero"
)

func handleInbox(log logger, fs *afero.Afero, inboxDir string, creds credentials.Credentials) error {
	log.Debug("Connecting to IMAP server")

	client, err := client.DialTLS(creds.IMAPServerAddress, nil)
	if err != nil {
		return fmt.Errorf("dialing: %w", err)
	}

	defer func() {
		_ = client.Logout()
	}()

	log.Debug("Logging in")

	if err = client.Login(creds.Username, creds.Password); err != nil {
		return fmt.Errorf("logging in: %w", err)
	}

	inbox, err := client.Select("INBOX", false)
	if err != nil {
		return fmt.Errorf("selecting INBOX: %w", err)
	}

	if inbox.Messages == 0 {
		return nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(inbox.Messages)

	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)

	done := make(chan error, 1)
	defer close(done)

	log.Debug("Initiating fetch")

	go func() {
		if err := client.Fetch(seqset, items, messages); err != nil {
			fmt.Println("Fetch error:", err)
		}
	}()

	go handleMessages(fs, inboxDir, section, messages, done)

	log.Debug("Waiting for message handling to finish")

	if err := <-done; err != nil {
		return fmt.Errorf("fetching messages: %w", err)
	}

	return nil
}

func handleMessages(fs *afero.Afero, inboxDir string, section imap.BodySectionName, messages chan *imap.Message, done chan error) {
	for {
		msg, ok := <-messages
		if !ok {
			break
		}

		parsedMessage, err := extractMessage(&section, msg)
		if err != nil {
			done <- fmt.Errorf("parsing message: %w", err)
		}

		err = fsconv.WriteMessageToDirectory(fs, inboxDir, parsedMessage)
		if err != nil {
			done <- fmt.Errorf("writing message to directory: %w", err)
		}
	}

	done <- nil
}

func extractMessage(section *imap.BodySectionName, rawMessage *imap.Message) (fsconv.Message, error) {
	resultMessage := fsconv.Message{}

	r := rawMessage.GetBody(section)
	if r == nil {
		fmt.Println("Server didn't returned message body")
		resultMessage.Body = strings.NewReader("<!-- no content -->")
	}

	mailReader, err := mail.CreateReader(r)
	if err != nil {
		return fsconv.Message{}, fmt.Errorf("creating mail reader: %w", err)
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
			return fsconv.Message{}, fmt.Errorf("reading mail part: %w", err)
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
