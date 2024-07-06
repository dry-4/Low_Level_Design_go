package main

import (
	"fmt"
	"math/rand"
	"publishSubscribe/pubsub"
	"time"
)

var availableTopics = map[string]string{
    "BTC": "BITCOIN",
    "ETH": "ETHEREUM",
    "DOT": "POLKADOT",
    "SOL": "SOLANA",
}


func pricePublisher(broker *pubsub.Broker) () {
	topicKeys := make([]string, 0, len(availableTopics))
   topicValues := make([]string, 0, len(availableTopics))

	for key, value := range availableTopics {
		topicKeys = append(topicKeys, key)
		_ = topicKeys // Ignore the result of append operation
		topicValues = append(topicValues, value)
	}

	for {
		randomValue := topicValues[rand.Intn(len(topicValues))]
		message := fmt.Sprintf("%f", rand.Float64())
		go broker.Publish(randomValue, message)
		fmt.Printf("Publishing %s to %s topic\n", message, randomValue)
		// go broker.Brodcast(message, topicValues)

		rand := 5
		time.Sleep(time.Duration(rand) * time.Second)
	}
}

func main() {
	broker := pubsub.NewBroker() // create a new broker
	s1 := broker.AddSubscriber() // add a new subscriber

	// subscribe BTC and ETH to s1.
   broker.Subscribe(s1, availableTopics["BTC"])
   broker.Subscribe(s1, availableTopics["ETH"])

	s2 := broker.AddSubscriber() // add a new subscriber

	// subscribe ETH and SOL to s2.
	broker.Subscribe(s2, availableTopics["ETH"])
	broker.Subscribe(s2, availableTopics["SOL"])

	// sleep for 5 sec, and then subscribe for topic DOT for s2
	go (func() {
		time.Sleep(3 * time.Second)
      broker.Subscribe(s2, availableTopics["DOT"])
	})();
	
	// sleep for 5 sec, and then unsubscribe for topic SOL for s2
	go (func(){
		time.Sleep(5*time.Second)
      broker.Unsubscribe(s2, availableTopics["SOL"])
      fmt.Printf("Total subscribers for topic ETH is %v\n", broker.GetSubscribers(availableTopics["ETH"]))
	})();
	
	// sleep for 10 sec, and then remove s2
	go (func ()  {
		time.Sleep(10*time.Second)
      broker.RemoveSubscriber(s2)
      fmt.Printf("Total subscribers for topic ETH is %v\n", broker.GetSubscribers(availableTopics["ETH"]))
	})();

	go pricePublisher(broker) // start the price publisher routine concurrently
	go s1.Listen() // start the listener for s1 concurrently
	go s2.Listen() // start the listener for s2 concurrently

	fmt.Scanln() // wait for a key press to exit, to prevent terminate of the program
	fmt.Println("Done!") 
}