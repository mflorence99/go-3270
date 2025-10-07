package device

// 🟧 Manage the device status by coordinating with the UI

func (device *Device) ResetStatus() {
	device.alarm = false
	device.cursorAt = 0
	device.erase = false
	device.error = false
	device.focussed = true
	device.locked = true
	device.message = ""
	device.numeric = false
	device.protected = false
	device.waiting = false
}

func (device *Device) SignalStatus() {
	status := map[string]any{
		"alarm":     device.alarm,
		"cursorAt":  device.cursorAt,
		"error":     device.error,
		"locked":    device.locked,
		"message":   device.message,
		"numeric":   device.numeric,
		"protected": device.protected,
		"waiting":   device.waiting,
	}
	device.SendGo3270Message(Go3270Message{eventType: "status", params: status})
	device.alarm = false
}

func (device *Device) StatusForAttributes(attrs *Attributes) {
	device.numeric = attrs.Numeric()
	device.protected = attrs.Protected()
}
