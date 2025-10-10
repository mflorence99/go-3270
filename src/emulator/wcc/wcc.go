package wcc

type WCC struct {
	Alarm    bool
	Reset    bool
	ResetMDT bool
	Unlock   bool
}

func New(char byte) *WCC {
	w := new(WCC)
	w.Alarm = (char & 0b00000100) != 0
	w.Reset = (char & 0b01000000) != 0
	w.ResetMDT = (char & 0b00000001) != 0
	w.Unlock = (char & 0b00000010) != 0
	return w
}

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

func (w *WCC) String() string {
	str := "WCC=[ "
	if w.Alarm {
		str += "ALARM "
	}
	if w.Reset {
		str += "RESET "
	}
	if w.ResetMDT {
		str += "-MDT "
	}
	if w.Unlock {
		str += "UNLOCK "
	}
	str += "]"
	return str
}
