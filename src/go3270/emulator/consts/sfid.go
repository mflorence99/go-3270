package consts

import (
	"fmt"
	"go3270/emulator/stream"
)

type SFID byte

type SFld struct {
	ID   SFID
	Info []byte
}

func (s SFld) String() string {
	return fmt.Sprintf("{ID %#02x, Info % #x}", byte(s.ID), s.Info)
}

const (
	QUERY_REPLY    SFID = 0x81
	READ_PARTITION SFID = 0x01
)

var sfids = map[SFID]string{
	0x81: "QUERY_REPLY",
	0x01: "READ_PARTITION",
}

// 🔥 eg: f3 0007 01 ff ff 03 00 81 ffef
func SFldsFromStream(out *stream.Outbound) []SFld {
	sflds := make([]SFld, 0)
	for out.HasNext() {
		len, _ := out.Next16()
		id, ok := out.Next()
		// 👇 there must be an ID
		if ok {
			// TODO 🔥 we can't account for this extra 0xFF!
			xtra, _ := out.Peek()
			if xtra == 0xFF {
				out.Skip(1)
			}
			var info []byte
			// 👇 a zero length can indicate the last field
			if len > 0 {
				info, _ = out.NextSlice(int(len) - 3)
			} else {
				info = out.Rest()
			}
			sfld := SFld{
				ID:   SFID(id),
				Info: info,
			}
			sflds = append(sflds, sfld)
		}
	}
	return sflds
}

func SFIDFor(s SFID) string {
	// 🔥 because we have not codified all of them, by a long shot!
	str, ok := sfids[s]
	if ok {
		return str
	} else {
		return fmt.Sprintf("%#02x", byte(s))
	}
}

func (s SFID) String() string {
	return SFIDFor(s)
}
