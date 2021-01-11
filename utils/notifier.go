package utils

import (
	"errors"
	"fmt"

	"github.com/kevinburke/twilio-go"
	"github.com/spf13/viper"
)

// Notifier allows for sending
// messages to a notification channel.
type Notifier struct {
	client    *twilio.Client
	sender    string
	recipient string
	nodeID    string
}

// NewNotifier ...
func NewNotifier(nodeID string) (*Notifier, error) {
	if len(viper.ConfigFileUsed()) == 0 {
		return nil, errors.New("config file at $HOME/.avalanchego/.avalanche-runner.yaml is missing")
	}

	accountSid := viper.GetString("twilio.accountSid")
	if len(accountSid) == 0 {
		return nil, errors.New("config file does not contain twilio.accountSid")
	}

	authToken := viper.GetString("twilio.authToken")
	if len(authToken) == 0 {
		return nil, errors.New("config file does not contain twilio.authToken")
	}

	sender := viper.GetString("twilio.sender")
	if len(sender) == 0 {
		return nil, errors.New("config file does not contain twilio.sender")
	}

	recipient := viper.GetString("twilio.recipient")
	if len(recipient) == 0 {
		return nil, errors.New("config file does not contain twilio.recipient")
	}

	return &Notifier{
		client:    twilio.NewClient(accountSid, authToken, nil),
		sender:    sender,
		recipient: recipient,
		nodeID:    nodeID,
	}, nil
}

// Info ...
func (n *Notifier) Info(message string) {
	fmt.Printf("NOTIFIER [INFO]: %s\n", message)
	_, err := n.client.Messages.SendMessage(
		n.sender,
		n.recipient,
		fmt.Sprintf("[INFO](%s): %s", n.nodeID, message),
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
		n.recipient,
		fmt.Sprintf("[ALERT](%s): %s", n.nodeID, message),
		nil,
	)
	if err != nil {
		fmt.Printf("NOTIFIER [ERROR]: %s\n", err.Error())
	}
}
