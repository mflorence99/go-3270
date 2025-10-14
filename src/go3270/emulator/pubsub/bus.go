package pubsub

import (
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

// 🟦 Type-safe publishers

func (b *Bus) Publish(topic string, args ...any) {
	handlers, ok := b.handlers[topic]
	if ok {
		for _, handler := range handlers {
			utils.Call(handler, args...)
		}
	}
}

func (b *Bus) PubClose() {
	b.Publish("close")
}

func (b *Bus) PubConfig(cfg Config) {
	b.Publish("config", cfg)
}

func (b *Bus) PubDump(dmp Dump) {
	b.Publish("dump", dmp)
}

func (b *Bus) PubFocus(focus bool) {
	b.Publish("focus", focus)
}

func (b *Bus) PubKeystroke(key Keystroke) {
	b.Publish("keystroke", key)
}

func (b *Bus) PubInbound(bytes []byte) {
	b.Publish("inbound", bytes)
}

func (b *Bus) PubOutbound(bytes []byte) {
	b.Publish("outbound", bytes)
}

func (b *Bus) PubPanic(msg string) {
	b.Publish("panic", msg)
}

func (b *Bus) PubStatus(stat Status) {
	b.Publish("status", stat)
}

// 🟦 Type-safe subscribers

func (b *Bus) Subscribe(topic string, fn interface{}) {
	// 🔥 ensure LIFO
	b.handlers[topic] = append([]interface{}{fn}, b.handlers[topic]...)
}

func (b *Bus) SubClose(fn func()) {
	b.Subscribe("close", fn)
}

func (b *Bus) SubConfig(fn func(cfg Config)) {
	b.Subscribe("config", fn)
}

func (b *Bus) SubDump(fn func(dmp Dump)) {
	b.Subscribe("config", fn)
}

func (b *Bus) SubFocus(fn func(focus bool)) {
	b.Subscribe("focus", fn)
}

func (b *Bus) SubKeystroke(fn func(key Keystroke)) {
	b.Subscribe("keystroke", fn)
}

func (b *Bus) SubInbound(fn func(bytes []byte)) {
	b.Subscribe("inbound", fn)
}

func (b *Bus) SubOutbound(fn func(bytes []byte)) {
	b.Subscribe("outbound", fn)
}

func (b *Bus) SubPanic(fn func(msg string)) {
	b.Subscribe("panic", fn)
}

func (b *Bus) SubStatus(fn func(stat Status)) {
	b.Subscribe("status", fn)
}

func (b *Bus) UnsubscribeAll() {
	b.handlers = make(map[string][]interface{})
}
