package sfld

import (
	"fmt"
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

// 🟧 3270 Structured Field ID

type SFld struct {
	ID   consts.SFID
	Info []byte
}

// 🟦 Constructor

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
				ID:   consts.SFID(id),
				Info: info,
			}
			sflds = append(sflds, sfld)
		}
	}
	return sflds
}

// 🟦 Stringer implementation

func (s SFld) String() string {
	return fmt.Sprintf("{ID %#02x, Info % #x}", byte(s.ID), s.Info)
}
