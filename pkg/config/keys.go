package config

const (
	// WorkingDirectory defines the current working directory, where the outbox, inbox, etc., are located.
	WorkingDirectory = "directory"
	// LogLevel defines the log level.
	LogLevel = "logLevel"

	// IMAPServerAddress defines the address of the IMAP server in a host:port format.
	IMAPServerAddress = "imapServerAddress"
	// SMTPServerAddress defines the address of the SMTP server in a host:port format.
	SMTPServerAddress = "smtpServerAddress"
)
