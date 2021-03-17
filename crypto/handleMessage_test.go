package crypto

import (
	"strings"
	"testing"
)

// Happy path.
func TestHandleMessage(t *testing.T) {
	expectedAddress := "0xDC25EF3F5B8A186998338A2ADA83795FBA2D695E"
	expectedText := "Sample text; with semicolons."
	message := expectedAddress + delimiter + expectedText
	expectedAddress = strings.ToLower(expectedAddress)

	address, text, err := HandleMessage(message)
	if err != nil {
		t.Errorf("HandleMessage() returned an error: %+v", err)
	}

	if expectedAddress != address {
		t.Errorf("HandleMessage() failed to return the expected address."+
			"\nexpected: %s\nreceived: %s", expectedAddress, address)
	}

	if expectedText != text {
		t.Errorf("HandleMessage() failed to return the expected text."+
			"\nexpected: %s\nreceived: %s", expectedText, text)
	}
}

// Error path: semicolon not included in message.
func TestHandleMessage_NoSemicolonError(t *testing.T) {
	expectedAddress := "0xDC25EF3F5B8A186998338A2ADA83795FBA2D695E"
	expectedText := "Sample text."
	message := expectedAddress + expectedText

	_, _, err := HandleMessage(message)
	if err == nil || !strings.Contains(err.Error(), "cannot process malformed message") {
		t.Errorf("HandleMessage() did not return an error when the message "+
			"lacked a semicolon: %+v", err)
	}
}
