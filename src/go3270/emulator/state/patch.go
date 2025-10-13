package state

type Patch struct {
	Alarm     *bool
	CursorAt  *int
	Error     *bool
	Locked    *bool
	Message   *string
	Numeric   *bool
	Protected *bool
	Waiting   *bool
}
