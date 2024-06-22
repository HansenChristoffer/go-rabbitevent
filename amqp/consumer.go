package amqp

import (
	"fmt"
	"log"

	"github.com/hansenchristoffer/go-rabbitevent/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer is responsible for managing the connection to the RabbitMQ server
// and consuming messages from specified queues.
type Consumer struct {
	connection *amqp.Connection
}

// NewConsumer creates a new instance of Consumer.
// It establishes a connection to the RabbitMQ server using the provided connection URL.
//
// Parameters:
// - connectionUrl: The URL to connect to the RabbitMQ server.
//
// Returns:
// - A pointer to the newly created Consumer instance.
// - An error if the connectionUrl is empty or if there is a failure in establishing the connection.
func NewConsumer(connectionUrl string) (*Consumer, error) {
	if len(connectionUrl) == 0 {
		return nil, fmt.Errorf("'connectionUrl' not allowed to be empty")
	}

	conn, err := amqp.Dial(connectionUrl)
	if err != nil {
		return nil, err
	}
	return &Consumer{connection: conn}, nil
}

// StartListening starts consuming messages from a specified queue and dispatches
// them to the appropriate listeners through the provided Dispatcher.
//
// Parameters:
// - queueName: The name of the queue to consume messages from. Must not be empty.
// - eventName: The name of the event to identify the consumer. Must not be empty.
// - dispatcher: A pointer to the Dispatcher instance that will handle the message dispatching.
// - messageType: An instance of the expected message type, used for unmarshalling the JSON message.
//
// Returns:
// - An error if the queueName or eventName is empty, or if there is a failure in setting up the consumer.
func (c *Consumer) StartListening(queueName string, eventName string,
	dispatcher *event.Dispatcher, messageType interface{}) error {

	if len(queueName) == 0 {
		return fmt.Errorf("'queueName' not allowed to be empty")
	}
	if len(eventName) == 0 {
		return fmt.Errorf("'eventName' not allowed to be empty")
	}

	ch, err := c.connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel -> %v", err)
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			log.Fatalf("error while closing amqp.Channel -> %v", err)
		}
	}(ch)

	messages, err := ch.Consume(
		queueName, // queue: The name of the queue to consume messages from.
		eventName, // consumer: The consumer identifier.
		true,      // auto-ack: If set to true, messages are acknowledged automatically.
		false,     // exclusive: If set to true, this consumer will be the only one consuming from the queue.
		false,     // no-local: If set to true, the server will not deliver messages to the consumer that were published on the same connection.
		false,     // no-wait: If set to true, the server will not respond to the method.
		nil,       // args: Additional arguments for the consumer.
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer -> %v", err)
	}

	fe := make(chan bool)

	go func() {
		for d := range messages {
			log.Printf("Received a message on %s[%s]: %s\n", eventName, queueName, d.Body)
			dispatcher.DispatchEvent(queueName, d.Body, messageType)
		}
	}()

	log.Printf(" Waiting for messages... Exit by pressing CTRL+C")
	<-fe
	return nil
}
