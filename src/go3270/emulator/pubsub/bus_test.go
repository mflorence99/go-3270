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
	bus.Subscribe(1, func(a, b int) { res = a + b })
	bus.Publish(1, 1, 2)
	assert.True(t, res == 3)
}

func Test_PubSub(t *testing.T) {
	bus := pubsub.NewBus()
	res := make([]string, 0)
	bus.Subscribe(1, func(x string) { res = append(res, x) })
	bus.Subscribe(1, func(x string) { res = append(res, x) })
	bus.Subscribe(1, func(x string) { res = append(res, x) })
	bus.Subscribe(2, func(x string) { res = append(res, x) })
	bus.Publish(1, "x")
	assert.True(t, slices.Equal(res, []string{"x", "x", "x"}))
}
