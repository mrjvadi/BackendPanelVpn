// Filename: events/event_bus.go
package events

import (
	"errors"
	"sync"
	"time"
	"github.com/google/uuid"
)

type Event struct {
	Name          string
	Payload       interface{}
	RequestID     string
	ResponseTopic string
}

type Handler func(Event)

type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewBus() *Bus {
	return &Bus{handlers: make(map[string][]Handler)}
}

func (b *Bus) Subscribe(name string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[name] = append(b.handlers[name], h)
}

func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	handlers := append([]Handler{}, b.handlers[e.Name]...)
	b.mu.RUnlock()
	for _, h := range handlers {
		go h(e)
	}
}

func (b *Bus) Request(topic string, payload interface{}, timeout time.Duration) (Event, error) {
	requestID := uuid.New().String()
	responseTopic := "response." + requestID
	responseChan := make(chan Event, 1)

	b.Subscribe(responseTopic, func(e Event) {
		responseChan <- e
	})
	
	b.Publish(Event{
		Name:          topic,
		Payload:       payload,
		RequestID:     requestID,
		ResponseTopic: responseTopic,
	})

	select {
	case response := <-responseChan:
		return response, nil
	case <-time.After(timeout):
		return Event{}, errors.New("request timed out")
	}
}