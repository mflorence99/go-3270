package device

// 🟧 Model 3270 WCC

type WCC struct {
	alarm    bool
	reset    bool
	resetMDT bool
	unlock   bool
}

func NewWCC(u8 byte) *WCC {
	wcc := new(WCC)
	wcc.alarm = (u8 & 0b00000100) != 0
	wcc.reset = (u8 & 0b01000000) != 0
	wcc.resetMDT = (u8 & 0b00000001) != 0
	wcc.unlock = (u8 & 0b00000010) != 0
	return wcc
}

func (wcc *WCC) DoAlarm() bool {
	return wcc.alarm
}

func (wcc *WCC) DoReset() bool {
	return wcc.reset
}

func (wcc *WCC) DoResetMDT() bool {
	return wcc.resetMDT
}

func (wcc *WCC) DoUnlock() bool {
	return wcc.unlock
}

func (wcc *WCC) ToByte() byte {
	var u8 byte = 0
	if wcc.alarm {
		u8 |= 0b00000100
	}
	if wcc.reset {
		u8 |= 0b01000000
	}
	if wcc.resetMDT {
		u8 |= 0b00000001
	}
	if wcc.unlock {
		u8 |= 0b00000010
	}
	return u8
}

func (wcc *WCC) ToString() string {
	str := "WCC=[ "
	if wcc.alarm {
		str += "ALARM "
	}
	if wcc.reset {
		str += "RESET "
	}
	if wcc.resetMDT {
		str += "-MDT "
	}
	if wcc.unlock {
		str += "UNLOCK "
	}
	str += "]"
	return str
}
