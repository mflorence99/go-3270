package consts

var (
	EUA byte = 0x12
	GE  byte = 0x08
	IC  byte = 0x13
	MF  byte = 0x2C
	PT  byte = 0x05
	RA  byte = 0x3C
	SA  byte = 0x28
	SBA byte = 0x11
	SF  byte = 0x1D
	SFE byte = 0x29
)

var orders = map[byte]string{
	0x05: "PT",
	0x08: "GE",
	0x11: "SBA",
	0x12: "EUA",
	0x13: "IC",
	0x1D: "SF",
	0x28: "SA",
	0x29: "SFE",
	0x2C: "MF",
	0x3C: "RA",
}

func OrderFor(order byte) string {
	return orders[order]
}
