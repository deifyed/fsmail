# fsmail

## Motivation

I want to handle my emails like I would handle my text files

## Usage

```shell
# Log in to your email provider. Credentials are stored in your secret store
fsmail login

# Synchronize your emails
fsmail sync

# Send an email by creating a file in ./outbox named after the recipient
cat <<EOF > important-email
---
From: me@example.com
To: lover@example.com
Subject: Missing you
---

Just wanted to let you know xoxo
EOF

# Then sync again to send the message
fsmail sync
```

## Installation

See [instructions](INSTALL.md)

## Configuration

N/A
