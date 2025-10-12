package bus

import (
	"fmt"
	"go3270/emulator/utils"
)

type Bus struct {
	handlers map[string][]interface{}
}

func NewBus() *Bus {
	b := new(Bus)
	b.handlers = make(map[string][]interface{})
	return b
}

func (b *Bus) Publish(topic string, args ...any) {
	_, ok := b.handlers[topic]
	if !ok {
		print(fmt.Sprintf("🔥 no handlers yet for %s", topic))
	}
	for ix := range b.handlers {
		utils.Call(b.handlers[ix], args)
	}
}

// 🔥 ensure LIFO
func (b *Bus) Subscribe(topic string, fn interface{}) {
	b.handlers[topic] = append([]interface{}{fn}, b.handlers[topic]...)
}

func (b *Bus) Unsubscribe(topic string) {
	delete(b.handlers, topic)
}

func (b *Bus) UnsubscribeAll() {
	b.handlers = make(map[string][]interface{})
}
