package types

import (
	"fmt"
)

// 🟧 3270 Structured Field ID

type SFID byte

// 🟦 Lookup tables

const (
	READ_PARTITION SFID = 0x01
	SET_REPLY_MODE SFID = 0x09
	QUERY_REPLY    SFID = 0x81
)

var sfids = map[SFID]string{
	0x01: "READ_PARTITION",
	0x09: "SET_REPLY_MODE",
	0x81: "QUERY_REPLY",
}

// 🟦 Stringer implementation

func SFIDFor(s SFID) string {
	// 🔥 because we have not codified all of them, by a long shot!
	if str, ok := sfids[s]; ok {
		return str
	} else {
		return fmt.Sprintf("%#02x", byte(s))
	}
}

func (s SFID) String() string {
	return SFIDFor(s)
}
