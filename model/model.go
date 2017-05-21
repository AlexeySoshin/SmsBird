package model

import (
	"fmt"
)

type Message struct {
	Recipient  int
	Originator string
	Message    string
	Total      int
	Count      int
	Reference  int
}

func (m *Message) String() string {
	return fmt.Sprintf("Recipient: %d, Originator: %s, Message: %s", m.Recipient, m.Originator, m.Message)
}
