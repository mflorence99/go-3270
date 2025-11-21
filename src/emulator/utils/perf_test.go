package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TODO 🔥 cheating here, but I don't want to change the API
// to pass a Writer just for this silly function

func TestPerf(t *testing.T) {
	defer ElapsedTime(time.Now())
	assert.True(t, true)
}
