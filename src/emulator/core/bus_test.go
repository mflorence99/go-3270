package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBus(t *testing.T) {
	bus := NewBus()
	bus.SubInitialize(func() {
		assert.Fail(t, "SubInit() should not be called here")
	})
	expected := []byte{0x01, 0x02, 0x03}
	var actual []byte
	bus.SubOutbound(func(chars []byte) {
		actual = chars
	})
	bus.PubOutbound(expected)
	assert.Equal(t, expected, actual, "smoke test passed")
}

func TestPubSubPanic(t *testing.T) {
	bus := NewBus()
	expected := "HELP!"
	var actual string
	bus.SubPanic(func(msg string) {
		actual = msg
	})
	bus.PubPanic(expected)
	assert.Equal(t, expected, actual, "Pub/SubPanic correct")
}

func TestPubSubRender(t *testing.T) {
	bus := NewBus()
	ok := false
	bus.SubRender(func() {
		ok = true
	})
	bus.PubRender()
	assert.True(t, ok, "Pub/SubRender correct")
}

func TestPubSubTrace(t *testing.T) {
	bus := NewBus()
	expected := "tick"
	var actual string
	bus.SubTrace(func(topic Topic, _ interface{}) {
		actual = topic.String()
	})
	bus.SubTick(func(counter int) {})
	bus.PubTick(0)
	assert.Equal(t, expected, actual, "SubTrace correct")
}

// TODO 🔥 we could call all of them, but ATM that's overkill
// mostly, if they compile they're correct after these smoke tests
