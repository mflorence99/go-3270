package device

// 🟧 Manage the device status by coordinating with the UI

func (device *Device) ResetStatus() {
	device.alarm = false
	device.cursorAt = 0
	device.erase = false
	device.error = false
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
	SendMessage(Message{bus: device.bus, eventType: "status", params: status})
}

func (device *Device) StatusForAttributes(attrs *Attributes) {
	device.numeric = attrs.IsNumeric()
	device.protected = attrs.IsProtected()
}
