# go-rabbitevent

go-rabbitevent is a Go library designed to facilitate the creation of event dispatchers on AMQP (RabbitMQ). It provides an easy-to-use interface to register event listeners, dispatch events to these listeners, and consume messages from RabbitMQ queues.

## Features

- Simplified registration and management of event listeners.
- Dynamic event dispatching based on message types.
- Easy integration with RabbitMQ using the `github.com/rabbitmq/amqp091-go` library.
- Support for concurrent message handling.

## Installation

First, install the library using `go get`:

```bash
go get github.com/hansenchristoffer/go-rabbitevent
```

Make sure to also install the github.com/rabbitmq/amqp091-go dependency:
```bash
go get github.com/rabbitmq/amqp091-go
```

## Usage
### Defining Custom Messages and Listeners
Define your custom message struct and implement the Listener interface for your custom listener:

```go
// CustomMessage represents the structure of the messages.
type CustomMessage struct {
    Field1 string `json:"field1"`
    Field2 int    `json:"field2"`
}

// CustomListener is an implementation of the event.Listener interface.
type CustomListener struct{}

// OnEvent processes the incoming CustomMessage.
func (l *CustomListener) OnEvent(message interface{}) {
    msg := message.(*CustomMessage)
    log.Printf("CustomListener received message: %+v", msg)
    // Add your message handling logic here
}
```

### Setting Up the Dispatcher and Consumer
Set up your event dispatcher and RabbitMQ consumer:

```go
func main() {
    // Create a new instance of ListenerRegistry.
    registry := event.NewListenerRegistry()

    // Register the CustomListener for the "custom_queue".
    registry.RegisterListener("custom_queue", &examples.CustomListener{})

    // Create a new EventDispatcher.
    dispatcher := event.NewEventDispatcher(registry)

    // Create a new RabbitMQ Consumer.
    consumer, err := amqp.NewConsumer("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("Failed to create RabbitMQ consumer: %v", err)
    }

    // Start listening on the "custom_queue" for messages of type CustomMessage.
    err = consumer.StartListening("custom_queue", "my_event", dispatcher, examples.CustomMessage{})
    if err != nil {
        log.Fatalf("Error while listening: %v", err)
    }

    // Keep the main function running indefinitely.
    select {}
}
```

### Detailed Explanation
**`amqp/consumer.go`**

The Consumer struct manages the connection to the RabbitMQ server and consumes messages from specified queues.

**`event/dispatcher.go`**
The Dispatcher struct is responsible for dispatching events to registered listeners.

**`event/registry.go`**
The ListenerRegistry struct manages the registration and retrieval of event listeners.

**`examples.go`**
The examples package demonstrates how to define custom messages and listeners, and how to set up the dispatcher and consumer.

## Contributing
Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgements
 - streadway/amqp - for the initial AMQP library for Go.
 - rabbitmq/amqp091-go - for the continued AMQP library for Go.