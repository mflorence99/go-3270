package wcc

type WCC struct {
	Alarm    bool
	Reset    bool
	ResetMDT bool
	Unlock   bool
}

func New(char byte) *WCC {
	wcc := new(WCC)
	wcc.Alarm = (char & 0b00000100) != 0
	wcc.Reset = (char & 0b01000000) != 0
	wcc.ResetMDT = (char & 0b00000001) != 0
	wcc.Unlock = (char & 0b00000010) != 0
	return wcc
}

func (wcc *WCC) Byte() byte {
	var u8 byte = 0
	if wcc.Alarm {
		u8 |= 0b00000100
	}
	if wcc.Reset {
		u8 |= 0b01000000
	}
	if wcc.ResetMDT {
		u8 |= 0b00000001
	}
	if wcc.Unlock {
		u8 |= 0b00000010
	}
	return u8
}

func (wcc *WCC) String() string {
	str := "WCC=[ "
	if wcc.Alarm {
		str += "ALARM "
	}
	if wcc.Reset {
		str += "RESET "
	}
	if wcc.ResetMDT {
		str += "-MDT "
	}
	if wcc.Unlock {
		str += "UNLOCK "
	}
	str += "]"
	return str
}
