package qr

import (
	"emulator/iface"
	"emulator/types"
	"encoding/binary"
)

// 🟧 Query Reply structured field

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

// 👁️ Query Reply {Alphanumeric Partitions) pp 6-23 to 6-25

type AlphanumericPartitions struct {
	SFID  types.SFID
	QCode types.QCode
	NA    byte
	M     uint16
	Flags byte
}

// 🟦 Constructor

func NewAlphanumericPartitions(cols, rows uint) AlphanumericPartitions {
	return AlphanumericPartitions{
		SFID:  types.QUERY_REPLY,
		QCode: types.ALPHANUMERIC_PARTITIONS,
		NA:    0x00,
		M:     uint16(cols * rows),
		Flags: 0x00,
	}
}

// 🟦 Public emitter function

func (s AlphanumericPartitions) Put(in iface.Inbound) {
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags and data
	chars = append(chars, s.NA)
	chars = binary.BigEndian.AppendUint16(chars, s.M)
	chars = append(chars, s.Flags)
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}
