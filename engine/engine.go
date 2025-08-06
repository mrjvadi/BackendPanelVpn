// Filename: engine/engine.go
package engine

import (
	"fmt"
	"github.com/mrjvadi/BackendPanelVpn/events"
	"github.com/mrjvadi/BackendPanelVpn/models"
	"time"
)

type Engine struct {
	bus *events.Bus
}

func New(bus *events.Bus) *Engine {
	return &Engine{bus: bus}
}

func (e *Engine) Start() {
	e.bus.Subscribe("service:job_request", e.handleServiceJob)
	e.bus.Subscribe("service:user_created", e.handleUserCreated)
}

func (e *Engine) handleServiceJob(evt events.Event) {
	fmt.Printf("[Engine] Received job (ID: %s), payload: %v. Simulating work...\n", evt.RequestID, evt.Payload)
	go func() {
		time.Sleep(2 * time.Second)
		responsePayload := fmt.Sprintf("Job '%v' completed successfully by Engine", evt.Payload)
		fmt.Printf("[Engine] Finished job (ID: %s). Sending callback.\n", evt.RequestID)
		e.bus.Publish(events.Event{
			Name:      evt.ResponseTopic,
			Payload:   responsePayload,
			RequestID: evt.RequestID,
		})
	}()
}

func (e *Engine) handleUserCreated(evt events.Event) {
	if user, ok := evt.Payload.(models.User); ok {
		fmt.Printf("[Engine] Noticed user creation: ID=%d, Username=%s\n", user.ID, user.Username)
	}
}
