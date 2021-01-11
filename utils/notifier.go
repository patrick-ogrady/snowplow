package utils

import (
	"fmt"

	"github.com/kevinburke/twilio-go"
)

// Notifier allows for sending
// messages to a notification channel.
type Notifier struct {
	client     *twilio.Client
	sender     string
	receipient string
}

// NewNotifier ...
func NewNotifier(
	accountSid string,
	authToken string,
	sender string,
	receipient string,
) *Notifier {
	return &Notifier{
		client:     twilio.NewClient(accountSid, authToken, nil),
		sender:     sender,
		receipient: receipient,
	}
}

// Info ...
func (n *Notifier) Info(message string) {
	fmt.Printf("NOTIFIER [INFO]: %s\n", message)
	_, err := n.client.Messages.SendMessage(
		n.sender,
		n.receipient,
		fmt.Sprintf("[INFO]: %s", message),
		nil,
	)
	if err != nil {
		fmt.Printf("NOTIFIER [ERROR]: %s\n", err.Error())
	}
}

// Alert ...
func (n *Notifier) Alert(message string) {
	fmt.Printf("NOTIFIER [ALERT]: %s\n", message)
	_, err := n.client.Messages.SendMessage(
		n.sender,
		n.receipient,
		fmt.Sprintf("[ALERT]: %s", message),
		nil,
	)
	if err != nil {
		fmt.Printf("NOTIFIER [ERROR]: %s\n", err.Error())
	}
}
