package types

// 🟧 3270 Orders

type Order byte

// 🟦 Lookup tables

const (
	EUA Order = 0x12
	GE  Order = 0x08
	IC  Order = 0x13
	MF  Order = 0x2C
	PT  Order = 0x05
	RA  Order = 0x3C
	SA  Order = 0x28
	SBA Order = 0x11
	SF  Order = 0x1D
	SFE Order = 0x29
)

var orders = map[Order]string{
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

// 🟦 Stringer implementation

func OrderFor(o Order) string {
	return orders[o]
}

func (o Order) String() string {
	return OrderFor(o)
}
