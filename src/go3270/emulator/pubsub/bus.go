package pubsub

import (
	"go3270/emulator/utils"
)

type Bus struct {
	handlers map[Topic][]interface{}
}

func NewBus() *Bus {
	b := new(Bus)
	b.handlers = make(map[Topic][]interface{})
	return b
}

func (b *Bus) Publish(topic Topic, args ...any) {
	handlers, ok := b.handlers[topic]
	if ok {
		for ix := range handlers {
			utils.Call(handlers[ix], args...)
		}
	}
}

// 🔥 ensure LIFO
func (b *Bus) Subscribe(topic Topic, fn interface{}) {
	b.handlers[topic] = append([]interface{}{fn}, b.handlers[topic]...)
}

func (b *Bus) Unsubscribe(topic Topic) {
	delete(b.handlers, topic)
}

func (b *Bus) UnsubscribeAll() {
	b.handlers = make(map[Topic][]interface{})
}
