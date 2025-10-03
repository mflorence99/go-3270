package device

import (
	"emulator/types"
	"emulator/utils"
	"time"
)

// 🟦 Subordinate to device.go: handle outbound data. Remember, "outbound" means from the application to the 3270 (this code)

func (device *Device) WriteCommands(out *OutboundDataStream) {
	defer utils.ElapsedTime(time.Now(), "WriteCommands")
	// 👇 dispatch on command
	switch device.command {
	case types.CommandLookup["RMA"]:
	case types.CommandLookup["EAU"]:
	case types.CommandLookup["EWA"]:
	case types.CommandLookup["W"]:
	case types.CommandLookup["RB"]:
	case types.CommandLookup["WSF"]:
	case types.CommandLookup["EW"]:
		device.InitializeBuffer()
		device.WriteOrdersAndData(out)
	case types.CommandLookup["RM"]:
	}
}

func (device *Device) WriteOrdersAndData(out *OutboundDataStream) {
	var lastAttrs *Attributes = NewAttributes([]uint8{0x00})
	for out.HasNext() {
		// 👇 look at each byte to see if it is an order
		order, _ := out.Next()
		// 👇 dispatch on order
		switch order {
		case types.OrderLookup["PT"]:
		case types.OrderLookup["GE"]:
		case types.OrderLookup["SBA"]:
			addr, _ := out.NextSlice(2)
			device.addr = utils.AddrFromBytes(addr)
			if device.addr >= device.size {
				SendMessage(Message{bus: device.bus, eventType: "panic", args: []any{"Data requires a device with a larger screen"}})
				return
			}
		case types.OrderLookup["EUA"]:
		case types.OrderLookup["IC"]:
			device.cursorAt = device.addr
		case types.OrderLookup["SF"]:
			attrs, _ := out.NextSlice(1)
			lastAttrs = NewAttributes(attrs)
			device.PutBuffer(0x00, lastAttrs)
		case types.OrderLookup["SA"]:
		case types.OrderLookup["SFE"]:
			count, _ := out.Next()
			attrs, _ := out.NextSlice(int(count) * 2)
			lastAttrs = NewAttributes(attrs)
			device.PutBuffer(0x00, lastAttrs)
		case types.OrderLookup["MF"]:
		case types.OrderLookup["RA"]:
		// 👇 if it isn't an order, it's data
		default:
			if order == 0x00 || order >= 0x40 {
				device.PutBuffer(utils.E2A([]uint8{order})[0], lastAttrs)
			}
		}
	}
}
