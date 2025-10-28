package buffer

type Fld []*Cell

func (f Fld) StartFld() (*Cell, bool) {
	if len(f) > 0 {
		return f[0], true
	}
	return nil, false
}
