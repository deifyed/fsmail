package fsconv

import (
	_ "embed"
	"io"
)

//go:embed file-template.md
var messageFileTemplate string

type Message struct {
	To      string
	From    string
	Cc      []string
	Subject string
	Body    io.Reader
}
