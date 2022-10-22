package sync

import (
	"fmt"
	"io"

	"github.com/deifyed/fsmail/pkg/credentials"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/spf13/afero"
)

func handleInbox(fs *afero.Afero, inboxDir string, creds credentials.Credentials) error {
	client, err := client.DialTLS(creds.IMAPServerAddress, nil)
	if err != nil {
		return fmt.Errorf("dialing: %w", err)
	}

	defer func() {
		_ = client.Logout()
	}()

	if err = client.Login(creds.Username, creds.Password); err != nil {
		return fmt.Errorf("logging in: %w", err)
	}

	mbox, err := client.Select("INBOX", false)
	if err != nil {
		return fmt.Errorf("selecting INBOX: %w", err)
	}

	if mbox.Messages == 0 {
		return nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(mbox.Messages)

	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	go func() {
		if err := client.Fetch(seqset, items, messages); err != nil {
			fmt.Println("Fetch error:", err)
		}
	}()

	msg := <-messages
	if msg == nil {
		fmt.Println("Server didn't returned message")
	}

	r := msg.GetBody(&section)
	if r == nil {
		fmt.Println("Server didn't returned message body")
	}

	mr, err := mail.CreateReader(r)
	if err != nil {
		return fmt.Errorf("creating mail reader: %w", err)
	}

	header := mr.Header

	if from, err := header.AddressList("From"); err == nil {
		fmt.Println("From:", from)
	}
	if to, err := header.AddressList("To"); err == nil {
		fmt.Println("To:", to)
	}
	if subject, err := header.Subject(); err == nil {
		fmt.Println("Subject:", subject)
	}

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("reading mail part: %w", err)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			body, err := io.ReadAll(p.Body)
			if err != nil {
				return fmt.Errorf("reading mail body: %w", err)
			}
			fmt.Println(string(body))
		case *mail.AttachmentHeader:
			filename, err := h.Filename()
			if err != nil {
				return fmt.Errorf("reading attachment filename: %w", err)
			}
			fmt.Println("Attachment:", filename)
		}
	}

	return nil
}
