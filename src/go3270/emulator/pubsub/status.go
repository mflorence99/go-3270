package pubsub

import (
	"fmt"
)

// 🟧 3270 status, as shared with Typescript UI

type Status struct {
	Alarm     bool
	CursorAt  int
	Error     bool
	Locked    bool
	Message   string
	Numeric   bool
	Protected bool
	Waiting   bool
}

// 🟦 Stringer implementation

func (s Status) String() string {
	str := fmt.Sprintf("CURSOR %d ", s.CursorAt)
	if s.Alarm {
		str += "ALARM "
	}
	if s.Locked {
		str += "LOCKED "
	}
	if s.Numeric {
		str += "NUM "
	}
	if s.Protected {
		str += "PROT "
	}
	if s.Waiting {
		str += "WAIT "
	}
	return fmt.Sprintf("%s%s", str, s.Message)
}
