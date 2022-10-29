package email

import (
	"fmt"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
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

	done := make(chan error, 1)
	defer close(done)

	log.Debug("Initiating fetch")

	go func() {
		if err := client.Fetch(seqset, items, messages); err != nil {
			fmt.Println("Fetch error:", err)
		}
	}()

	convertedMessages := make([]Message, 0)
	go handleMessages(section, messages, done, convertedMessages)

	log.Debug("Waiting for message handling to finish")

	if err := <-done; err != nil {
		return nil, fmt.Errorf("fetching messages: %w", err)
	}

	return convertedMessages, nil
}
