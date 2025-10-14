package pubsub

import (
	"fmt"
)

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

func (s Status) String() string {
	str := ""
	if s.Alarm {
		str += "ALARM "
	}
	if s.Locked {
		str += "LOCK "
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
