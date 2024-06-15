package event

import "sync"

// Listener defines the interface that all event listeners must implement.
type Listener interface {
	OnEvent(message interface{})
}

// ListenerRegistry manages the registration and retrieval of event listeners.
type ListenerRegistry struct {
	listeners map[string][]Listener // A map of queue names to their respective listeners.
	mu        sync.RWMutex          // A RWMutex to ensure thread-safe access to the listeners map.
}

// NewListenerRegistry creates a new instance of ListenerRegistry.
// It initializes the internal listeners map.
func NewListenerRegistry() *ListenerRegistry {
	return &ListenerRegistry{
		listeners: make(map[string][]Listener),
	}
}

// RegisterListener registers an event listener for a specific queue.
// The function ensures thread-safe modification of the listeners map.
//
// Parameters:
// - queueName: The name of the queue for which the listener is being registered.
// - listener: The listener to register for the specified queue.
func (r *ListenerRegistry) RegisterListener(queueName string, listener Listener) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listeners[queueName] = append(r.listeners[queueName], listener)
}

// GetListeners retrieves all listeners registered for a specific queue.
// The function ensures thread-safe read access to the listeners map.
//
// Parameters:
// - queueName: The name of the queue whose listeners are to be retrieved.
//
// Returns:
// A slice of Listener instances registered for the specified queue.
// If no listeners are registered for the queue, an empty slice is returned.
func (r *ListenerRegistry) GetListeners(queueName string) []Listener {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.listeners[queueName]
}
