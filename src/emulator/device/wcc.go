package device

// 🟧 Model 3270 WCC

type WCC struct {
	alarm    bool
	reset    bool
	resetMDT bool
	unlock   bool
}

func NewWCC(byte uint8) *WCC {
	wcc := new(WCC)
	wcc.alarm = (byte & 0b00000100) != 0
	wcc.reset = (byte & 0b01000000) != 0
	wcc.resetMDT = (byte & 0b00000001) != 0
	wcc.unlock = (byte & 0b00000010) != 0
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

func (wcc *WCC) ToByte() uint8 {
	var byte uint8 = 0
	if wcc.alarm {
		byte |= 0b00000100
	}
	if wcc.reset {
		byte |= 0b01000000
	}
	if wcc.resetMDT {
		byte |= 0b00000001
	}
	if wcc.unlock {
		byte |= 0b00000010
	}
	return byte
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
