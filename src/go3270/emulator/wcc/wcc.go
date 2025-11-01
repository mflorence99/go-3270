package wcc

// 🟧 3270 WCC (write control character)

type WCC struct {
	Alarm    bool
	Reset    bool
	ResetMDT bool
	Unlock   bool
}

// 🟦 Constructor

func NewWCC(char byte) WCC {
	return WCC{
		Alarm:    (char & 0b00000100) != 0,
		Reset:    (char & 0b01000000) != 0,
		ResetMDT: (char & 0b00000001) != 0,
		Unlock:   (char & 0b00000010) != 0,
	}
}

// 🟦 Public functions

func (w *WCC) Byte() byte {
	var u8 byte = 0
	if w.Alarm {
		u8 |= 0b00000100
	}
	if w.Reset {
		u8 |= 0b01000000
	}
	if w.ResetMDT {
		u8 |= 0b00000001
	}
	if w.Unlock {
		u8 |= 0b00000010
	}
	return u8
}
