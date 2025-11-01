package consts

// 🟧 Query Reply / List responses

type QCode byte

// 🟦 Lookup tables

const (
	ALPHANUMERIC_PARTITIONS QCode = 0x84
	CHARACTER_SETS          QCode = 0x85
	COLOR_SUPPORT           QCode = 0x86
	DDM                     QCode = 0x95
	FIELD_OUTLINING         QCode = 0x8C
	FIELD_VALIDATION        QCode = 0x8A
	HIGHLIGHTING            QCode = 0x87
	IMPLICIT_PARTITION      QCode = 0xA6
	REPLY_MODES             QCode = 0x88
	RPQ_NAMES               QCode = 0xA1
	SUMMARY                 QCode = 0x80
	USABLE_AREA             QCode = 0x81
)

var qcodes = map[QCode]string{
	0x80: "SUMMARY",
	0x81: "USABLE_AREA",
	0x84: "ALPHANUMERIC_PARTITIONS",
	0x85: "CHARACTER_SETS",
	0x86: "COLOR_SUPPORT",
	0x87: "HIGHLIGHTING",
	0x88: "REPLY_MODES",
	0x8A: "FIELD_VALIDATION",
	0x8C: "FIELD_OUTLINING",
	0x95: "DDM",
	0xA1: "RPQ_NAMES",
	0xA6: "IMPLICIT_PARTITION",
}

// 🟦 Stringer implementation

func QCodeFor(s QCode) string {
	return qcodes[s]
}

func (s QCode) String() string {
	return QCodeFor(s)
}
