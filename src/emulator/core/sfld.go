package core

import (
	"emulator/types"
	"fmt"
)

// 🟧 3270 Structured Field ID

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

type SFld struct {
	ID   types.SFID
	Info []byte
}

// 🟦 Constructor

// 👁️ Introduction  pp 5-4 to 5-5
func SFldsFromStream(out *Outbound) []SFld {
	sflds := make([]SFld, 0)
	for out.HasNext() {
		len := out.MustNext16()
		// 👇 there must be an ID
		if id, ok := out.Next(); ok {
			// TODO 🔥 we can't account for this extra 0xfF!
			xtra := out.MustPeek()
			if xtra == 0xfF {
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
