package device

// 🟧 Send messages to the UI for action

// 👁️ go3270.go subscribes to them on the bus, then uses dispatchEvent to route them to mediator.ts in the UI. All this is necessary so that go test-able code can communicate with the ui without using syscall/js

type Go3270Message struct {
	args      []any
	u8s       []byte
	eventType string
	params    map[string]any
}

func (device *Device) SendGo3270Message(msg Go3270Message) {
	device.bus.Publish("go3270", msg.eventType, msg.u8s, msg.params, msg.args)
}
