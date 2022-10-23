package fsconv

import (
	"fmt"
	"path"
	"testing"

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
					Recipient: "me@example.com",
					Subject:   "This is a test subject",
					Body:      "Mock content",
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
				err := messageToFile(fs, workDir, message)
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