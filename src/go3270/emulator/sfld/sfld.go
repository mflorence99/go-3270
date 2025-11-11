package sfld

import (
	"fmt"
	"go3270/emulator/stream"
	"go3270/emulator/types"
)

// 🟧 3270 Structured Field ID

type SFld struct {
	ID   types.SFID
	Info []byte
}

// 🟦 Constructor

func SFldsFromStream(out *stream.Outbound) []SFld {
	sflds := make([]SFld, 0)
	for out.HasNext() {
		len := out.MustNext16()
		id, ok := out.Next()
		// 👇 there must be an ID
		if ok {
			// TODO 🔥 we can't account for this extra 0xFF!
			xtra := out.MustPeek()
			if xtra == 0xFF {
				out.MustSkip(1)
			}
			var info []byte
			// 👇 a zero length can indicate the last field
			if len > 0 {
				info = out.MustNextSlice(int(len) - 3)
			} else {
				info = out.Rest()
			}
			sfld := SFld{
				ID:   types.SFID(id),
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
