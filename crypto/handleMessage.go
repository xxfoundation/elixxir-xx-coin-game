package crypto

import (
	"github.com/pkg/errors"
	"strings"
)

const delimiter = ";"

// HandleMessage splits the message at the semicolon into an address and text.
// The address is returned in lowercase. An error is only returned when the
// semicolon is not present.
func HandleMessage(message string) (string, string, error) {
	// Split message into two parts at first semicolon
	msgParts := strings.SplitN(message, delimiter, 2)

	// Return an error if the semicolon was not present
	if len(msgParts) < 2 {
		return "", "", errors.New("cannot process malformed message; format " +
			"is: ethAddress;message")
	}

	return strings.ToLower(msgParts[0]), msgParts[1], nil
}
