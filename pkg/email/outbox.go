package email

import (
	"fmt"
	"strconv"
	"strings"
)

func parseServerAddress(serverAddress string) (string, int, error) {
	parts := strings.Split(serverAddress, ":")

	host := parts[0]

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("converting port from string to int: %w", err)
	}

	return host, port, nil
}
