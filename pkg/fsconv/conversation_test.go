package fsconv

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"
	"testing"
	"text/template"

	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestMessageToFile(t *testing.T) {
	testCases := []struct {
		name                string
		withMessages        []Message
		expectExistingFiles []string
	}{
		{
			name: "Should generate expected file with one message",
			withMessages: []Message{
				{
					To:      "me@example.com",
					Subject: "This is a test subject",
					Body:    strings.NewReader("Mock content"),
				},
			},
			expectExistingFiles: []string{"/work/This-is-a-test-subject"},
		},
		{
			name: "Should generate expected file with one message and multiple recipients",
			withMessages: []Message{
				{
					To:      "me@example.com",
					Cc:      []string{"someone@example.com", "else@example.com"},
					Subject: "This is a test subject",
					Body:    strings.NewReader("Mock content"),
				},
			},
			expectExistingFiles: []string{"/work/This-is-a-test-subject"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := &afero.Afero{Fs: afero.NewMemMapFs()}
			workDir := "/work"

			for _, message := range tc.withMessages {
				err := WriteMessageToDirectory(fs, workDir, message)
				assert.NoError(t, err)
			}

			g := goldie.New(t)

			for _, file := range tc.expectExistingFiles {
				exists, err := fs.Exists(file)
				assert.NoError(t, err)

				assert.True(t, exists, "Expected file to exist: %s", file)

				content, err := fs.ReadFile(file)
				assert.NoError(t, err)

				g.Assert(t, fmt.Sprintf("%s-%s", t.Name(), path.Base(file)), content)
			}
		})
	}
}

type testFile struct {
	filepath string
	content  io.Reader
}

const fileTemplate = `---
To: {{ .To }}
From: {{ .From }}
Subject: {{ .Subject }}
---

{{ .Body }}`

func createTestFileContent(t *testing.T, from, to, subject, body string) io.Reader {
	templ := template.Must(template.New("").Parse(fileTemplate))

	buf := bytes.Buffer{}

	err := templ.Execute(&buf, struct {
		From    string
		To      string
		Subject string
		Body    string
	}{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    body,
	})
	assert.NoError(t, err)

	return &buf
}

func TestDirectoryToMessages(t *testing.T) {
	testCases := []struct {
		name           string
		withFiles      []testFile
		expectMessages []Message
	}{
		{
			name: "Should find and extract correctly a single message",
			withFiles: []testFile{
				{
					filepath: "/This-is-a-test-subject",
					content:  createTestFileContent(t, "me@example.com", "you@example.com", "mock subject", "mock body"),
				},
			},
			expectMessages: []Message{
				{
					From:    "me@example.com",
					To:      "you@example.com",
					Subject: "mock subject",
					Body:    bytes.NewBuffer([]byte("mock body")),
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := &afero.Afero{Fs: afero.NewMemMapFs()}

			for _, file := range tc.withFiles {
				err := fs.WriteReader(file.filepath, file.content)
				assert.NoError(t, err)
			}

			messages, err := DirectoryToMessages(fs, "/")
			assert.NoError(t, err)

			actualMessageMap := messagesAsMap(messages)

			for _, message := range tc.expectMessages {
				assert.Contains(t, actualMessageMap, message.Subject)

				actualMessage := actualMessageMap[message.Subject]

				assert.Equal(t, message.From, actualMessage.From)
				assert.Equal(t, message.To, actualMessage.To)
				assert.Equal(t, message.Subject, actualMessage.Subject)

				actualBody, err := io.ReadAll(actualMessage.Body)
				assert.NoError(t, err)

				expectedBody, err := io.ReadAll(message.Body)
				assert.NoError(t, err)

				assert.Equal(t, string(expectedBody), string(actualBody))
			}

			assert.Equal(t, tc.expectMessages, messages)
		})
	}
}

func messagesAsMap(messages []Message) map[string]Message {
	m := make(map[string]Message)

	for _, message := range messages {
		m[message.Subject] = message
	}

	return m
}
