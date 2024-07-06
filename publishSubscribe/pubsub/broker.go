package pubsub

import (
	"fmt"
	"sync"
)

type Subscribers map[string]*Subscriber

// Broker - represents a broker in a pub/sub system.
type Broker struct {
	subscribers Subscribers // map of subscribers id:subscriber
	topics map[string]Subscribers // map of topics to subscribers
	mutex sync.RWMutex // lock
}


// NewBroker - creates a new broker.
func NewBroker() *Broker {
	return &Broker{
		subscribers: Subscribers{},
		topics: map[string]Subscribers{},
	}
}

// Subscribe - subscribes a subscriber to a topic.
func (broker *Broker) Subscribe(subscriber *Subscriber, topic string) {
	broker.mutex.RLock()
	defer broker.mutex.RUnlock()
	if broker.topics[topic] == nil {
		broker.topics[topic] = Subscribers{}
	}
	subscriber.AddTopic(topic)
	broker.topics[topic][subscriber.id] = subscriber
	fmt.Printf("%s Subscribed for topic: %s\n", subscriber.id, topic)
}

// unsubscribe - unsubscribes a subscriber from a topic.
func (broker *Broker) Unsubscribe(subscriber *Subscriber, topic string) {
	broker.mutex.RLock()
    defer broker.mutex.RUnlock()
    delete(broker.topics[topic], subscriber.id)
    subscriber.RemoveTopic(topic)
    fmt.Printf("%s Unsubscribed for topic: %s\n", subscriber.id, topic)
}

// AddSubscriber - adds a subscriber to the broker.
func (broker *Broker) AddSubscriber() (*Subscriber) {
	broker.mutex.Lock()
   defer broker.mutex.Unlock()
	id, subscriber := CreateNewSubscriber()
	broker.subscribers[id] = subscriber
	return subscriber
}

// RemoveSubscriber - removes a subscriber from the broker.
func (broker *Broker) RemoveSubscriber(subscriber *Subscriber) {
	// Unsubscribe the subscriber from all topics.
	for topic := range subscriber.topics {
		broker.Unsubscribe(subscriber, topic)
	}

	broker.mutex.Lock()
	delete(broker.subscribers, subscriber.id)
	broker.mutex.Unlock()
	subscriber.Deactivate()
	fmt.Printf("Removed subscriber %s\n", subscriber.id)
}

// GetSubscribers - returns the number of subscribers for a topic.
func (broker *Broker) GetSubscribers(topic string) int {
	broker.mutex.RLock()
	defer broker.mutex.RUnlock()
	return len(broker.topics[topic])
}

// Broadcast - broadcasts a message to all subscribers of a topic.
func (broker *Broker) Broadcast(msg string, topics []string) {
	broker.mutex.RLock()
	defer broker.mutex.RUnlock()

	for _, topic := range topics {
		for _, subscriber := range broker.topics[topic] {
			if !subscriber.IsActive() {
				continue
			}

			message := NewMessage(topic, msg)
			go (func(subscriber *Subscriber) {
				subscriber.Signal(message)
			})(subscriber)
		}
	}
}


// Publish - publishes a message to a topic.
func (broker *Broker) Publish(topic string, msg string) {
	
	broker.mutex.RLock()
	brokerTopics := broker.topics[topic]
	broker.mutex.RUnlock()
	
	for _, subscriber := range brokerTopics {
		if !subscriber.IsActive() {
			return
		}
		message := NewMessage(topic, msg)
		
		go (func(subscriber *Subscriber) {
			subscriber.Signal(message)
		})(subscriber)
		fmt.Printf("Published & signaled %s to %s topic\n", msg, topic)
	}
}

func DisplayBrok() {
	fmt.Println("Hello, World! from custom.go broker.go")
}