// Copyright (c) 2021 patrick-ogrady
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package notifier

import (
	"errors"
	"fmt"

	"github.com/kevinburke/twilio-go"
	"github.com/spf13/viper"
)

const (
	info   = "INFO"
	alert  = "ALERT"
	status = "STATUS"
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
		return nil, errors.New(
			"config file at $HOME/.avalanchego/.snowplow.yaml is missing",
		)
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

func (n *Notifier) sendMessage(kind string, message string) {
	if n == nil {
		return
	}

	_, err := n.client.Messages.SendMessage(
		n.sender,
		n.recipient,
		fmt.Sprintf("[%s](%s): %s", kind, n.nodeID, message),
		nil,
	)
	if err != nil {
		fmt.Printf("notifier error: %s\n", err.Error())
	}
}

// Info ...
func (n *Notifier) Info(message string) {
	n.sendMessage(info, message)
}

// Alert ...
func (n *Notifier) Alert(message string) {
	n.sendMessage(alert, message)
}

// Status ...
func (n *Notifier) Status(message string) {
	n.sendMessage(status, message)
}
