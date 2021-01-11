package utils

// Notifier allows for sending
// messages to a notification channel.
type Notifier struct {
}

// NewNotifier ...
func NewNotifier() *Notifier {
	return &Notifier{}
}

// Info ...
func (n *Notifier) Info(message string) {
}

// Alert ...
func (n *Notifier) Alert(message string) {
}
