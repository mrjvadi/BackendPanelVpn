package events

import "testing"

func TestBusPublishSubscribe(t *testing.T) {
	bus := NewBus()
	called := false
	bus.Subscribe("test:event", func(Event) { called = true })
	bus.Publish(Event{Name: "test:event"})
	if !called {
		t.Fatalf("handler was not called")
	}
}
