package convert

import (
	"bytes"
	"html/template"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToMessage(t *testing.T) {
	testCases := []struct {
		name          string
		withContent   io.Reader
		expectMessage Message
	}{
		{
			name:        "Should successfully convert a simple message",
			withContent: newEmailContent(t, "me@example.com", "you@example.com", "testing", "such long mock body"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m, err := ToMessage(tc.withContent)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectMessage, m)
		})
	}
}

func TestToReader(t *testing.T) {
	testCases := []struct {
		name          string
		with          Message
		expectContent string
	}{
		{
			name: "Should convert a plain message",
			with: Message{
				From:    "me@example.com",
				To:      "you@example.com",
				Subject: "testing",
				Body:    "such long mock body",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			expectedContent := newEmailContent(t, tc.with.From, tc.with.To, tc.with.Subject, tc.with.Body)
			got := ToReader(tc.with)

			expectedContentAsBytes, err := io.ReadAll(expectedContent)
			assert.NoError(t, err)

			gotContentAsBytes, err := io.ReadAll(got)
			assert.NoError(t, err)

			assert.Equal(t, string(expectedContentAsBytes), string(gotContentAsBytes))
		})
	}
}

const emailTemplate = `---
To: {{ .To }}
From: {{ .From }}
Subject: {{ .Subject }}
---

{{ .Body }}
`

func newEmailContent(t *testing.T, from, to, subject, body string) io.Reader {
	t.Helper()

	tem, err := template.New("email").Parse(emailTemplate)
	assert.NoError(t, err)

	var b bytes.Buffer

	err = tem.Execute(&b, struct {
		To      string
		From    string
		Subject string
		Body    string
	}{
		To:      to,
		From:    from,
		Subject: subject,
		Body:    body,
	})
	assert.NoError(t, err)

	return &b
}
