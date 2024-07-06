package pubsub

import (
	"fmt"
)

// Message is a struct that represents a message in a pub/sub system.
type Message struct {
	// Topic is the topic string that the message is published to.
	topic string
	// Body is the body of the message.
	body  string
}


// NewMessage - creates a new Message with the given topic and body.
func NewMessage(topic, body string) *Message {
	return &Message{topic: topic, body: body}
}

// GetBody - returns the body of the message.
func (msg *Message) GetTopic() string {
	return msg.topic
}

// GetBody - returns the body of the message.
func (msg *Message) GetMessageBody() string {
	return msg.body
}

func DisplayMsg() {
	fmt.Println("Hello, World! from custom.go message.go")
}

