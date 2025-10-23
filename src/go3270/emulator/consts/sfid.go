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

func FromStream(out *stream.Outbound) []SFld {
	sflds := make([]SFld, 0)
	for out.HasNext() {
		len, _ := out.Next16()
		id, _ := out.Next()
		if len > 0 {
			info, _ := out.NextSlice(int(len) - 3)
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
