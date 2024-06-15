package event

import (
	"encoding/json"
	"log"
	"reflect"
)

// Dispatcher is responsible for dispatching events to registered listeners.
type Dispatcher struct {
	registry *ListenerRegistry
}

// NewEventDispatcher creates a new instance of Dispatcher.
// It initializes the dispatcher with a given ListenerRegistry.
//
// Parameters:
// - registry: A pointer to an instance of ListenerRegistry that manages event listeners.
//
// Returns:
// A pointer to the newly created Dispatcher instance.
func NewEventDispatcher(registry *ListenerRegistry) *Dispatcher {
	return &Dispatcher{registry: registry}
}

// DispatchEvent dispatches an event to all listeners registered for a specific queue.
// It dynamically constructs the event message from JSON, ensuring that it is passed
// to listeners as the correct type.
//
// Parameters:
// - queueName: The name of the queue for which the event is being dispatched.
// - message: The raw JSON message to be dispatched to the listeners.
// - messageType: An instance of the expected message type, used for unmarshalling the JSON message.
//
// The function logs an error if it fails to unmarshal the message. It retrieves the list of listeners
// for the given queue from the registry and dispatches the event to each listener in a separate goroutine
// to ensure non-blocking, concurrent handling.
func (d *Dispatcher) DispatchEvent(queueName string, message []byte, messageType interface{}) {
	// Create a new instance of the messageType to hold the unmarshalled message.
	event := reflect.New(reflect.TypeOf(messageType)).Interface()

	// Unmarshal the JSON message into the event instance.
	err := json.Unmarshal(message, event)
	if err != nil {
		log.Printf("Error unmarshaling message -> %v", err)
		return
	}

	// Retrieve the list of listeners for the specified queue.
	listeners := d.registry.GetListeners(queueName)

	// Dispatch the event to each listener in a separate goroutine.
	for _, listener := range listeners {
		go listener.OnEvent(event)
	}
}
