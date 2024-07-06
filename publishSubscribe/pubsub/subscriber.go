package pubsub

import (
	"crypto/rand"
	"fmt"
	"log"
	"sync"
)


type Subscriber struct {
	 id string // id of subscriber
    messages chan* Message // messages channel
    topics map[string]bool // topics it is subscribed to.
    active bool // if given subscriber is active
    mutex sync.RWMutex // lock
}

// NewSubscriber - creates a new subscriber with the given id.
func CreateNewSubscriber() (string, *Subscriber) {
	b := make([]byte, 16)
	_, err := rand.Read(b)

	if err != nil {
		log.Fatal(err)
	}

	id := fmt.Sprintf("%X", b)

	return id, &Subscriber{
		id: id,
		messages: make(chan *Message),
		topics: map[string]bool{},
		active: true,
	}
}


// AddTopic - adds a topic to the subscriber.
func (s * Subscriber) AddTopic(topic string)(){
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.topics[topic] = true
}

// RemoveTopic - removes a topic from the subscriber.
func (s * Subscriber) RemoveTopic(topic string)(){
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	delete(s.topics, topic)
}

// GetTopics - returns all topics the subscriber is subscribed to.
func (s * Subscriber) GetTopics() ([]string){
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	topics := []string{}
	for topic := range s.topics {
		topics = append(topics, topic)
	}
	return topics
}

// GetID - returns the id of the subscriber.
func (s * Subscriber) GetID() (string){
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.id
}

// IsActive - returns if the subscriber is active.
func (s * Subscriber) IsActive() (bool){
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.active
}

// Deactivate - deactivates the subscriber.
func (s * Subscriber) Deactivate()(){
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.active = false
}

// Activate - activates the subscriber.
func (s * Subscriber) Activate()(){
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.active = true
}

// Signal - signals the subscriber to listen for messages.
func (s * Subscriber) Signal(msg *Message) {
	s.mutex.RLock()
   defer s.mutex.RUnlock()
   if s.active{
      s.messages <- msg
   }
}


// Listen - listens for messages on the subscriber's message channel.
func (s * Subscriber) Listen() {
	for {
   	if msg, ok := <- s.messages; ok {
      	fmt.Printf("Subscriber %s, received: %s from topic: %s\n", s.id, msg.GetMessageBody(), msg.GetTopic())
      }
    }
}

func DisplaySubs() {
	fmt.Println("Hello, World! from custom.go subsciber.go")
}