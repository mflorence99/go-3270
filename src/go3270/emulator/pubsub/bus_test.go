package pubsub_test

import (
	"go3270/emulator/pubsub"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Smoke(t *testing.T) {
	bus := pubsub.NewBus()
	res := 0
	bus.Subscribe("x", func(a, b int) { res = a + b })
	bus.Publish("x", 1, 2)
	assert.True(t, res == 3)
}

func Test_PubSub(t *testing.T) {
	bus := pubsub.NewBus()
	res := make([]string, 0)
	bus.Subscribe("x", func(x string) { res = append(res, x) })
	bus.Subscribe("x", func(x string) { res = append(res, x) })
	bus.Subscribe("x", func(x string) { res = append(res, x) })
	bus.Subscribe("y", func(x string) { res = append(res, x) })
	bus.Publish("x", "x")
	assert.True(t, slices.Equal(res, []string{"x", "x", "x"}))
}
