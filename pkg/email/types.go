package email

import "io"

type Credentials struct {
	IMAPServerAddress string
	Username          string
	Password          string
}

type Message struct {
	From    string
	To      string
	Subject string
	Body    io.Reader
}

type logger interface {
	Debug(string)
	Debugf(string, ...interface{})
}
