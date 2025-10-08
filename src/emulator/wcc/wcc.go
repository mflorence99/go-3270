package wcc

// 🟧 Model 3270 WCC

type WCC struct {
	alarm    bool
	reset    bool
	resetMDT bool
	unlock   bool
}

func New(ch byte) *WCC {
	wcc := new(WCC)
	wcc.alarm = (ch & 0b00000100) != 0
	wcc.reset = (ch & 0b01000000) != 0
	wcc.resetMDT = (ch & 0b00000001) != 0
	wcc.unlock = (ch & 0b00000010) != 0
	return wcc
}

func (wcc *WCC) Alarm() bool {
	return wcc.alarm
}

func (wcc *WCC) Reset() bool {
	return wcc.reset
}

func (wcc *WCC) ResetMDT() bool {
	return wcc.resetMDT
}

func (wcc *WCC) Unlock() bool {
	return wcc.unlock
}

func (wcc *WCC) Byte() byte {
	var u8 byte = 0
	if wcc.Alarm() {
		u8 |= 0b00000100
	}
	if wcc.Reset() {
		u8 |= 0b01000000
	}
	if wcc.ResetMDT() {
		u8 |= 0b00000001
	}
	if wcc.Unlock() {
		u8 |= 0b00000010
	}
	return u8
}

func (wcc *WCC) String() string {
	str := "WCC=[ "
	if wcc.Alarm() {
		str += "ALARM "
	}
	if wcc.Reset() {
		str += "RESET "
	}
	if wcc.ResetMDT() {
		str += "-MDT "
	}
	if wcc.Unlock() {
		str += "UNLOCK "
	}
	str += "]"
	return str
}
