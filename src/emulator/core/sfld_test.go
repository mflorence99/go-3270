package core

import (
	"emulator/types"
	"emulator/utils"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSFldsFromStream(t *testing.T) {
	raw := []any{
		[]byte{0x00, 0x06},
		types.READ_PARTITION,
		[]byte{0x01, 0x02, 0x03},
		[]byte{0x00, 0x05},
		types.QUERY_REPLY,
		[]byte{0x04, 0x05},
	}
	bus := NewBus()
	out := NewOutbound(utils.Flatten2Bytes(raw, strings.ToUpper), bus)
	sflds := SFldsFromStream(out)
	assert.Equal(t, 2, len(sflds), "two SFlds found")
	assert.Equal(t, types.READ_PARTITION, sflds[0].ID, "READ_PARTITION is first")
	assert.Equal(t, types.QUERY_REPLY, sflds[1].ID, "QUERY_REPLY is next")
}
