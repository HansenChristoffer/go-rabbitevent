package examples

import (
	"log"

	"github.com/hansenchristoffer/go-rabbitevent/amqp"
	"github.com/hansenchristoffer/go-rabbitevent/event"
)

// CustomMessage represents the structure of the messages that will be handled
// by the CustomListener. Modify this struct to fit the schema of your specific messages.
type CustomMessage struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

// CustomListener is an implementation of the event.Listener interface.
// It defines the logic to handle events of type CustomMessage.
type CustomListener struct{}

// OnEvent is called when an event of type CustomMessage is received.
// This method processes the incoming message.
//
// Parameters:
// - message: The incoming message to be processed. It is expected to be of type CustomMessage.
func (l *CustomListener) OnEvent(message interface{}) {
	msg := message.(*CustomMessage)
	log.Printf("CustomListener received message: %+v", msg)
	// Add your message handling logic here
}

func main() {
	// Create a new instance of ListenerRegistry to manage event listeners.
	registry := event.NewListenerRegistry()

	// Register a CustomListener for the "custom_queue".
	registry.RegisterListener("custom_queue", &CustomListener{})

	// Create a new EventDispatcher to dispatch events to the registered listeners.
	dispatcher := event.NewEventDispatcher(registry)

	// Create a new RabbitMQ Consumer using the provided connection URL.
	consumer, err := amqp.NewConsumer("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ consumer: %v", err)
	}

	// Start listening on the "custom_queue" for messages of type CustomMessage.
	// The dispatcher will handle dispatching messages to the appropriate listeners.
	err = consumer.StartListening("custom_queue", "my_event", dispatcher, CustomMessage{})
	if err != nil {
		log.Fatalf("got error while listening -> %v", err)
	}

	// Keep the main function running indefinitely to allow continuous message processing.
	select {}
}
