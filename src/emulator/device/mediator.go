package device

import (
	"github.com/asaskevich/EventBus"
)

// 🟧 Send messages to the UI for action

// 👁️ go3270.go subscribes to them on the bus, then uses dispatchEvent to route them to mediator.ts in the UI. All this is necessary so that go test-able code can communicate with the ui without using syscall/js

type Message struct {
	args      []any
	bus       EventBus.Bus
	bytes     []byte
	eventType string
	params    map[string]any
}

func SendMessage(msg Message) {
	msg.bus.Publish("go3270", msg.eventType, msg.bytes, msg.params, msg.args)
}
