package types

// 🟧 3270 status, as shared with Typescript UI

type Status struct {
	Alarm     bool
	CursorAt  uint
	Error     bool
	Insert    bool
	Locked    bool
	Message   string
	Numeric   bool
	Protected bool
	Waiting   bool
}

type Patch struct {
	Alarm     *bool
	CursorAt  *uint
	Error     *bool
	Insert    *bool
	Locked    *bool
	Message   *string
	Numeric   *bool
	Protected *bool
	Waiting   *bool
}
