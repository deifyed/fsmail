# fssmtp

## Motivation

I want to handle my emails like I would handle my text files

## Usage

```shell
# Log in to your email provider. Credentials are stored in your secret store
fssmtp login

# Synchronize your emails
fssmtp sync

# Send an email by creating a file in ./outbox named after the recipient
echo "Hey!\n\nI miss you" > outbox/lover@example.com

# Then sync again to send the message
fssmtp sync
```

## Installation

See [instructions](INSTALL.md)

## Configuration

N/A
