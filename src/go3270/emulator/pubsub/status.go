package pubsub

import (
	"fmt"
	"strings"
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
	var b strings.Builder
	fmt.Fprintf(&b, "CURSOR %d ", s.CursorAt)
	if s.Alarm {
		b.WriteString("ALARM ")
	}
	if s.Locked {
		b.WriteString("LOCKED ")
	}
	if s.Numeric {
		b.WriteString("NUM ")
	}
	if s.Protected {
		b.WriteString("PROT ")
	}
	if s.Waiting {
		b.WriteString("WAIT ")
	}
	return fmt.Sprintf("%s%s", b.String(), s.Message)
}
